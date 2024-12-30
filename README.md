The concept of this project is to increase throughput for event driven systems by creating concurrent operations across two services using websockets. 
One service would be the handling tasks and the other service would be handling the initial request and persistence.
The /process endpoint is upgraded to a websocket using gorilla websocket
Read more about it hear https://github.com/gorilla/websocket
