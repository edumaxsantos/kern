# Kern

This project aims to be a simple way to communicate with your arduino projects using a terminal app.

In this case, the software for the arduino was also implemented using go (TinyGo), but you could use any other language that you can think of, just have to implement the communication protocol.

As of the time of writing this readme, the cli can do the following:

- Connect to the user selected port
- Communicate to the user selected pin
- Send values (0/1) to the selected pin
- Get the current value of the selected pin
- Change the pin

For the future, I'm thinking about allowing the user to have multiple pins controlled at the same time, so the user does not have to change it every time.

Another possible improvement is to allow the user to send analog value to analog pins. Right now it only allow for digital/binary value.
