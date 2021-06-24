donut.go is a Go implementation of [Andy Sloan's donut.c](http://www.a1k0n.net/2011/07/20/donut-math.html).

It demonstrates:

* 3d rendering in a terminal window with ascii characters
* determining the terminal screen size
* using a channel as a timer

#### This is the code in C:
```
             k;double sin()
         ,cos();main(){float A=
       0,B=0,i,j,z[1760];char b[
     1760];printf("\x1b[2J");for(;;
  ){memset(b,32,1760);memset(z,0,7040)
  ;for(j=0;6.28>j;j+=0.07)for(i=0;6.28
 >i;i+=0.02){float c=sin(i),d=cos(j),e=
 sin(A),f=sin(j),g=cos(A),h=d+2,D=1/(c*
 h*e+f*g+5),l=cos      (i),m=cos(B),n=s\
in(B),t=c*h*g-f*        e;int x=40+30*D*
(l*h*m-t*n),y=            12+15*D*(l*h*n
+t*m),o=x+80*y,          N=8*((f*e-c*d*g
 )*m-c*d*e-f*g-l        *d*n);if(22>y&&
 y>0&&x>0&&80>x&&D>z[o]){z[o]=D;;;b[o]=
 ".,-~:;=!*#$@"[N>0?N:0];}}/*#****!!-*/
  printf("\x1b[H");for(k=0;1761>k;k++)
   putchar(k%80?b[k]:10);A+=0.04;B+=
     0.02;}}/*****####*******!!=;:~
       ~::==!!!**********!!!==::-
         .,~~;;;========;;;:~-.
             ..,--------,*/
```

#### This is the result in Go:

		                            
		        @@@@@$$             
		      @@@$$#####$$          
		    $@$$$##*!*!**#$         
		   #$$$$#*!===!!!*##        
		   $$$$#*!=;:::;=!*##       
		  #$$$$#*=;~-,-~:=**##      
		  #$$$#**=:-  .-:=*###*     
		 *#$$$##!=~    .~=*##$#     
		 *#$$$##*=~     -;*#$$#=    
		 *######*!:     -=*#$$#!    
		 !*#####**=     ~!#$$$#*    
		 !*######**=    ;#$$$$#!    
		 =**#######*!  =#$@@@$#!    
		 =!**###########$@@@@$#!    
		 :!!**######$$$$$@@@$#*;    
		  =!!**#####$$$$$$$$$#!;    
		  :=!****####$$$$$$##*=     
		   :=!!!***########**=:     
		   -;;=!!!**********!;,     
		    ,:;==!!!*!!!!!!!;-      
		     ,~;;====!====;:-       
		       -~:;;;;;;:~~         
		         .---~~-,           
                
At its core, it’s a framebuffer and a Z-buffer into which I render pixels. Since it’s just rendering relatively low-resolution ASCII art, I massively cheat. All it does is plot pixels along the surface of the torus at fixed-angle increments, and does it densely enough that the final result looks solid. The “pixels” it plots are ASCII characters corresponding to the illumination value of the surface at each point: .,-~:;=!*#$@ from dimmest to brightest. No raytracing required.

So how do we do that? Well, let’s start with the basic math behind 3D perspective rendering. The following diagram is a side view of a person sitting in front of a screen, viewing a 3D object behind it.

![](https://www.a1k0n.net/img/perspective.png)
To render a 3D object onto a 2D screen, we project each point (x,y,z) in 3D-space onto a plane located z’ units away from the viewer, so that the corresponding 2D position is (x’,y’). Since we’re looking from the side, we can only see the y and z axes, but the math works the same for the x axis (just pretend this is a top view instead). This projection is really easy to obtain: notice that the origin, the y-axis, and point (x,y,z) form a right triangle, and a similar right triangle is formed with (x’,y’,z’). Thus the relative proportions are maintained:
![](https://latex.codecogs.com/svg.image?\frac{y'}{z'}=\frac{y}{z})
![](https://latex.codecogs.com/svg.image?y'=\frac{yz'}{z'})
So to project a 3D coordinate to 2D, we scale a coordinate by the screen distance z’. Since z’ is a fixed constant, and not functionally a coordinate, let’s rename it to K1, so our projection equation becomes ![](https://latex.codecogs.com/svg.image?(x',y')=(\frac{K_1x}{z},\frac{K_1y}{z})).

. We can choose K1 arbitrarily based on the field of view we want to show in our 2D window. For example, if we have a 100x100 window of pixels, then the view is centered at (50,50); and if we want to see an object which is 10 units wide in our 3D space, set back 5 units from the viewer, then K1 should be chosen so that the projection of the point x=10, z=5 is still on the screen with x’ < 50: 10K1/5 < 50, or K1 < 25.

When we’re plotting a bunch of points, we might end up plotting different points at the same (x’,y’) location but at different depths, so we maintain a z-buffer which stores the z coordinate of everything we draw. If we need to plot a location, we first check to see whether we’re plotting in front of what’s there already. It also helps to compute z-1 =1z
and use that when depth buffering because:

* z-1 = 0 corresponds to infinite depth, so we can pre-initialize our z-buffer to 0 and have the background be infinitely far away 
* we can re-use z-1 when computing x’ and y’: Dividing once and multiplying by z-1 twice is cheaper than dividing by z twice.

Now, how do we draw a donut, AKA torus? Well, a torus is a solid of revolution, so one way to do it is to draw a 2D circle around some point in 3D space, and then rotate it around the central axis of the torus. Here is a cross-section through the center of a torus:
![](https://www.a1k0n.net/img/torusxsec.png)
So we have a circle of radius R1 centered at point (R2,0,0), drawn on the xy-plane. We can draw this by sweeping an angle — let’s call it θ — from 0 to 2π:
![](https://latex.codecogs.com/svg.image?(x,y,z)=(R_2,0,0)+(R_1\cos\theta,R_1\sin\theta,0)).

Now we take that circle and rotate it around the y-axis by another angle — let’s call it φ. To rotate an arbitrary 3D point around one of the cardinal axes, the standard technique is to multiply by a [rotation matrix](https://en.wikipedia.org/wiki/Rotation_matrix).

We also want the whole donut to spin around on at least two more axes for the animation. They were called A and B in the original code: it was a rotation about the x-axis by A and a rotation about the z-axis by B.

## Usage ##

To run:

	$ go run donut.go
