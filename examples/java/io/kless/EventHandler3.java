package io.kless;

import java.nio.ByteBuffer;

class EventHandler3 implements EventHandlerInterface {

    public Response eventHandler(Context context, Request req) {
        System.out.println("Inside event handler...");

        Response response = new Response();
        response.setBody(ByteBuffer.wrap("Hello world!".getBytes()));
        
        return response;
    }

}
