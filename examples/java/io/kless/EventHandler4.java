package io.kless;

import java.nio.ByteBuffer;

class EventHandler1 implements EventHandlerInterface {

    public Response eventHandler(Context context, Request req) {
        System.out.println("Inside event handler...");

        System.out.println("Request:");
        System.out.println(new String(req.getBody().array()));

        Response response = new Response();
        response.setBody(ByteBuffer.wrap("Request received".getBytes()));
        
        return response;
    }

}
