
import io.swagger.client.*;
import io.swagger.client.auth.*;
import io.swagger.client.model.*;
import io.swagger.client.api.DefaultApi;

import java.io.File;
import java.util.*;

public class DefaultApiExample {

    public static void main(String[] args) {
        ApiClient defaultClient = Configuration.getDefaultApiClient();


        DefaultApi apiInstance = new DefaultApi();
        UpdateStatusRequest body = new UpdateStatusRequest(); // UpdateStatusRequest | 
        try {
            StatusResponse result = apiInstance.statusPost(body);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling DefaultApi#statusPost");
            e.printStackTrace();
        }
    }
}