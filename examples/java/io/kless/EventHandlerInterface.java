package io.kless;

interface EventHandlerInterface {
    Response eventHandler(Context context, Request req);
}