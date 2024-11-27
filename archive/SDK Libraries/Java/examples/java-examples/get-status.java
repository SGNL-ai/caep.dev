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
        String streamId = "streamId_example"; // String | 
        try {
            StatusResponse result = apiInstance.statusGet(streamId);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling DefaultApi#statusGet");
            e.printStackTrace();
        }
    }
}