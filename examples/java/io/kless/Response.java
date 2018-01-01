package io.kless;

import java.util.Map;
import java.nio.ByteBuffer;

class Response {
    private Map<String,String> headers = null;
    private ByteBuffer body = null;

    public Map<String,String> getHeaders() {
        return headers;
    }

    public void setHeaders(Map<String,String> m) {
        this.headers = m;
    }

    public ByteBuffer getBody() {
        return body;
    }

    public void setBody(ByteBuffer buf) {
        this.body = buf;
    }

}