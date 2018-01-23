package io.kless;

import java.sql.*;

class EventHandler1 implements EventHandlerInterface {

    public Response eventHandler(Context context, Request req) {
        System.out.println("Inside event handler...");

        try {
            String url = "jdbc:postgresql://192.168.1.233/test?user=test&password=test";
            Connection conn = DriverManager.getConnection(url);

            // Do something useful here...

            conn.close();
        } catch (Exception e) {
            e.printStackTrace();
        }

        System.out.println("Leaving event handler...");

        return null;
    }

}


