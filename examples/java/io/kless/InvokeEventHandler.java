package io.kless;

import java.io.BufferedInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.Channels;
import java.nio.channels.WritableByteChannel;
import java.util.Map;
import java.util.HashMap;
import java.util.List;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import com.sun.net.httpserver.Headers;

class InvokeEventHandler {

    public static void main(String[] args) {

        String eventHandlerClass = System.getProperty("klessEventHandlerClass");
        String eventHandlerPortNumber = System.getProperty("klessEventHandlerPortNumber");

        System.out.println("Event handler class = " + eventHandlerClass);
        System.out.println("Event handler port number = " + eventHandlerPortNumber);

        int portNumber = Integer.parseInt(eventHandlerPortNumber);

        System.out.println("Starting...");

        try {
            HttpServer server = HttpServer.create(new InetSocketAddress("0.0.0.0", portNumber), 0);
            server.createContext("/", new HttpHandler() {
                @Override
                public void handle(HttpExchange t) throws IOException {
                    Context context = new Context();

                    Request req = new Request();
                    Response resp = null;

                    Map<String,String> requestHeaders = new HashMap<String,String>();
                    Headers headers = t.getRequestHeaders();
                    for (String headerName : headers.keySet()) {
                        List<String> headerValues = headers.get(headerName);
                        requestHeaders.put(headerName, headerValues.get(0));
                    }
                    req.setHeaders(requestHeaders);

                    ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
                    BufferedInputStream inputStream = new BufferedInputStream(t.getRequestBody());
                    byte[] buf = new byte[4096];
                    int numRead = 0;
                    while ((numRead = inputStream.read(buf, 0, 4096)) > 0) {
                        byteArrayOutputStream.write(buf, 0, numRead);
                    }
                    req.setBody(ByteBuffer.wrap(byteArrayOutputStream.toByteArray()));

                    EventHandlerInterface eventHandler;
                    try {
                        eventHandler = (EventHandlerInterface) Class.forName(eventHandlerClass).newInstance();

                        resp = eventHandler.eventHandler(context, req);
                    } catch (InstantiationException e) {
                        System.err.println("Unable to instantiate event handler from class = " + eventHandlerClass + ". The class must implement the io.kless.EventHandlerInterface interface.");
                    } catch (ClassNotFoundException e) {
                        System.err.println("Event handler class = " + eventHandlerClass + " was not found. Please validate the class name.");
                    } catch (Exception e) {
                        System.err.println("Exception caught from event handler class = " + eventHandlerClass);
                        e.printStackTrace();
                    }
                    
                    if (null != resp) {
                        if (null != resp.getHeaders()) {
                            Headers httpResponseHeaders = t.getResponseHeaders();
                            Map<String,String> kubelessResponseHeaders = resp.getHeaders();
                            for (String headerName : kubelessResponseHeaders.keySet()) {
                                httpResponseHeaders.add(headerName, kubelessResponseHeaders.get(headerName));
                            }
                        }
                        if (null != resp.getBody()) {
                            t.sendResponseHeaders(200, resp.getBody().limit());
                            OutputStream outputStream = t.getResponseBody();
                            WritableByteChannel channel = Channels.newChannel(outputStream);
                            channel.write(resp.getBody());
                            channel.close();
                            outputStream.close();
                        } else {
                            t.sendResponseHeaders(200, 0);
                        }
                    } else {
                        t.sendResponseHeaders(200, 0);
                    }
                    t.close();
                }
            });
            server.setExecutor(null);
            server.start();
        } catch (Exception e) {
            e.printStackTrace();
        }  
    }

}


