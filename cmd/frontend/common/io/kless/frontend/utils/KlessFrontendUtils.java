package io.kless.frontend.utils;

import java.net.InetAddress;
import java.net.URL;
import java.net.HttpURLConnection;
import java.io.OutputStream;
import java.io.IOException;
import java.util.Map;
import java.util.Date;
import java.util.concurrent.TimeUnit;

import org.influxdb.InfluxDBFactory;
import org.influxdb.InfluxDB;
import org.influxdb.dto.BatchPoints;
import org.influxdb.dto.Point;

/**
 * Utilities that can be used from Java-based frontends to perform common tasks.
 */
public class KlessFrontendUtils {

    /**
     * Send a message to the event handler container in this pod.
     * 
     * @param headers Headers to forward to the event handler
     * @params body Body to forward to the event handler
     * @params podName The name of this pod
     * @param eventHandler Name of the event handler
     * @param version Version of the event handler
     * @param namespace Namespace of the event handler
     */
    public void sendMessage(Map<String,Object> headers, String body, String podName, String eventHandler, String version, String namespace) throws Exception {
        long startTime = System.currentTimeMillis();

        HttpURLConnection connection = (HttpURLConnection) new URL("http://localhost:8080").openConnection();

        connection.setRequestMethod("POST");
        connection.setDoOutput(true);
        if (null != headers) {
            for (String key : headers.keySet()) {
                connection.setRequestProperty(key, headers.get(key).toString());
            }
        }

        connection.connect();

        OutputStream out = connection.getOutputStream();
        byte[] buf = body.getBytes("UTF-8");
        out.write(buf, 0, buf.length);
        out.close();

        System.out.println("Response code = " + connection.getResponseCode());
        int contentLength = connection.getContentLength();

        connection.disconnect();

        long respTime = System.currentTimeMillis() - startTime;

        System.out.println("Sending event to influxdb");
        sendEventToInfluxDB(eventHandler, namespace, version, podName, buf.length, contentLength, respTime);
        System.out.println("Event sent");
    }

    /**
     * Store the info about this event that was processed in the InfluxDB database.
     *
     * @param eventHandler Name of the event handler
     * @param namespace Namespace of the event handler
     * @param version Version of the event handler
     * @param reqSize The size of the request body in bytes
     * @param respSize The size of the response body in bytes
     * @param respTime The time it took the event handler to process this event in ms
     */
    private void sendEventToInfluxDB(String eventHandler, String namespace, String version, String podName, int reqSize, int respSize, long respTime) {
        try {
            InfluxDB influxDB = InfluxDBFactory.connect("http://kless-influxdb.kless:8086");
            String dbName = "klessdb";
            influxDB.createDatabase(dbName);

            BatchPoints batchPoints = BatchPoints
                .database(dbName)
                .tag("namespace", namespace) 
                .tag("handler", eventHandler)
                .tag("version", version)
                .tag("podname", podName)
                .retentionPolicy("autogen")
                .build();
            Point point1 = Point.measurement("kless")
                    .time(System.currentTimeMillis(), TimeUnit.MILLISECONDS)
                    .addField("reqSize", reqSize)
                    .addField("respSize", respSize)
                    .addField("respTime", respTime)
                    .build();
            batchPoints.point(point1);

            influxDB.write(batchPoints);
        } catch (Exception e) {
            System.err.println("Unable to send event to influxdb");
            e.printStackTrace();
        }
    }

}