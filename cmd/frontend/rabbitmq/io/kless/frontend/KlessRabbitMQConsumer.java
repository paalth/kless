package io.kless.frontend;

import java.net.InetAddress;
import java.util.Map;
import java.util.Date;
import java.io.IOException;

import com.rabbitmq.client.ConnectionFactory;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Consumer;
import com.rabbitmq.client.DefaultConsumer;
import com.rabbitmq.client.Envelope;
import com.rabbitmq.client.AMQP;

import io.kless.frontend.utils.KlessFrontendUtils;

class KlessRabbitMQConsumer {

    public static final void main(String[] args) throws Exception {
        System.out.println("Starting RabbitMQ frontend...");

        String hostname = System.getenv("KLESS_FRONTEND_HOSTNAME");
        String queueName = System.getenv("KLESS_FRONTEND_QUEUE_NAME");
        String consumerType = System.getenv("KLESS_FRONTEND_TYPE");
        String username = System.getenv("KLESS_FRONTEND_USERNAME");
        String password = System.getenv("KLESS_FRONTEND_PASSWORD");
        String consumerTag = InetAddress.getLocalHost().getHostName();
        String eventHandler = System.getenv("KLESS_EVENT_HANDLER_NAME");
        String namespace = System.getenv("KLESS_EVENT_HANDLER_NAMESPACE");
        String version = System.getenv("KLESS_EVENT_HANDLER_VERSION");

        new KlessRabbitMQConsumer().startFrontend(hostname, queueName, consumerTag, username, password, eventHandler, version, namespace);
    }

    private void startFrontend(String hostname, String queueName, String consumerTag, String username, String password, String eventHandler, String version, String namespace) throws Exception {
        System.out.println("Connecting, hostname = " + hostname + ", queue = " + queueName + ", consumer = " + consumerTag + ", username = " + username);

        ConnectionFactory factory = new ConnectionFactory();
        factory.setHost(hostname);
        factory.setUsername(username);
        factory.setPassword(password);
        Connection connection = factory.newConnection();
        Channel channel = connection.createChannel();

        channel.queueDeclare(queueName, true, false, false, null);

        channel.basicQos(1);

        System.out.println("Connected. Receiving message...");

        Consumer consumer = new DefaultConsumer(channel) {
            @Override
            public void handleDelivery(String consumerTag, Envelope envelope, AMQP.BasicProperties properties, byte[] body) throws IOException {
                String messageId = properties.getMessageId();
                Date timestamp = properties.getTimestamp();
                System.out.println("Message received: msg id = " + messageId + ", timestamp = " + timestamp);

                Map<String, Object> headers = properties.getHeaders();
                if (null != headers) {
                    System.out.println("Message headers:");
                    for (String key : headers.keySet()) {
                        System.out.println(key + ":" + headers.get(key));
                    }
                }

                String message = new String(body, "UTF-8");
                System.out.println("Message body: " + message);

                try {
                    new KlessFrontendUtils().sendMessage(headers, message, consumerTag, eventHandler, version, namespace);

                    channel.basicAck(envelope.getDeliveryTag(), false);
                } catch (Exception e) {
                    System.err.println("Unable to send message to event handler, not ack'ing message");
                    e.printStackTrace();
                }
            }
        };

        boolean autoAck = false;
        channel.basicConsume(queueName, autoAck, consumerTag, consumer);

        while (true) {
            Thread.sleep(30000);
            System.out.println("Receiving messages...");
        }
    }

}