package io.kless;

import java.util.Map;

class EventHandler2 implements EventHandlerInterface {

    public Response eventHandler(Context context, Request req) {
        System.out.println("Inside event handler example 2...");

        System.out.println("Request headers:");
        Map<String,String> headers = req.getHeaders();
        for (String headerName : headers.keySet()) {
            System.out.println(headerName + " = " + headers.get(headerName));
        }
        
        return null;
    }

}


