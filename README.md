donut.go is a Go implementation of [Andy Sloan's donut.c](http://www.a1k0n.net/2011/07/20/donut-math.html).

It demonstrates:

* 3d rendering in a terminal window with ascii characters
* determining the terminal screen size
* using a channel as a timer

## Screenshot ##

		                            
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
                        

           
## Usage ##

To run:

	$ go run donut.go
