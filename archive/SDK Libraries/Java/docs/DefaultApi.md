# DefaultApi

All URIs are relative to *https://caep-ssf.sbx.sgnl.host/ssf*

Method | HTTP request | Description
------------- | ------------- | -------------
[**statusGet**](DefaultApi.md#statusGet) | **GET** /status | 
[**statusPost**](DefaultApi.md#statusPost) | **POST** /status | 
[**streamsPollPost**](DefaultApi.md#streamsPollPost) | **POST** /streams/poll | 
[**streamsPost**](DefaultApi.md#streamsPost) | **POST** /streams | 
[**subjectsaddPost**](DefaultApi.md#subjectsaddPost) | **POST** /subjects:add | 
[**subjectsremovePost**](DefaultApi.md#subjectsremovePost) | **POST** /subjects:remove | 

<a name="statusGet"></a>
# **statusGet**
> StatusResponse statusGet(streamId)



Get the status of the stream with CAEP transmitter.

### Example
```java
// Import classes:
//import io.swagger.client.ApiClient;
//import io.swagger.client.ApiException;
//import io.swagger.client.Configuration;
//import io.swagger.client.auth.*;
//import io.swagger.client.api.DefaultApi;

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
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **streamId** | **String**|  | [optional]

### Return type

[**StatusResponse**](StatusResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

<a name="statusPost"></a>
# **statusPost**
> StatusResponse statusPost(body)



Update the stream status to the CAEP transmitter.

### Example
```java
// Import classes:
//import io.swagger.client.ApiClient;
//import io.swagger.client.ApiException;
//import io.swagger.client.Configuration;
//import io.swagger.client.auth.*;
//import io.swagger.client.api.DefaultApi;

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
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**UpdateStatusRequest**](UpdateStatusRequest.md)|  |

### Return type

[**StatusResponse**](StatusResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="streamsPollPost"></a>
# **streamsPollPost**
> PollEventsResponse streamsPollPost(body)



Poll and/or acknowledge events from the stream.

### Example
```java
// Import classes:
//import io.swagger.client.ApiClient;
//import io.swagger.client.ApiException;
//import io.swagger.client.Configuration;
//import io.swagger.client.auth.*;
//import io.swagger.client.api.DefaultApi;

ApiClient defaultClient = Configuration.getDefaultApiClient();


DefaultApi apiInstance = new DefaultApi();
PollEventsRequest body = new PollEventsRequest(); // PollEventsRequest | 
try {
    PollEventsResponse result = apiInstance.streamsPollPost(body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling DefaultApi#streamsPollPost");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**PollEventsRequest**](PollEventsRequest.md)|  |

### Return type

[**PollEventsResponse**](PollEventsResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="streamsPost"></a>
# **streamsPost**
> CreateStreamResponse streamsPost(body)



create a stream with CAEP transmitter.

### Example
```java
// Import classes:
//import io.swagger.client.ApiClient;
//import io.swagger.client.ApiException;
//import io.swagger.client.Configuration;
//import io.swagger.client.auth.*;
//import io.swagger.client.api.DefaultApi;

ApiClient defaultClient = Configuration.getDefaultApiClient();


DefaultApi apiInstance = new DefaultApi();
CreateStreamRequest body = new CreateStreamRequest(); // CreateStreamRequest | 
try {
    CreateStreamResponse result = apiInstance.streamsPost(body);
    System.out.println(result);
} catch (ApiException e) {
    System.err.println("Exception when calling DefaultApi#streamsPost");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CreateStreamRequest**](CreateStreamRequest.md)|  |

### Return type

[**CreateStreamResponse**](CreateStreamResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="subjectsaddPost"></a>
# **subjectsaddPost**
> subjectsaddPost(body)



Add a subject to be listened by the stream.

### Example
```java
// Import classes:
//import io.swagger.client.ApiClient;
//import io.swagger.client.ApiException;
//import io.swagger.client.Configuration;
//import io.swagger.client.auth.*;
//import io.swagger.client.api.DefaultApi;

ApiClient defaultClient = Configuration.getDefaultApiClient();


DefaultApi apiInstance = new DefaultApi();
AddSubjectRequest body = new AddSubjectRequest(); // AddSubjectRequest | 
try {
    apiInstance.subjectsaddPost(body);
} catch (ApiException e) {
    System.err.println("Exception when calling DefaultApi#subjectsaddPost");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AddSubjectRequest**](AddSubjectRequest.md)|  |

### Return type

null (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

<a name="subjectsremovePost"></a>
# **subjectsremovePost**
> subjectsremovePost(body)



Remove a listened subject from the stream.

### Example
```java
// Import classes:
//import io.swagger.client.ApiClient;
//import io.swagger.client.ApiException;
//import io.swagger.client.Configuration;
//import io.swagger.client.auth.*;
//import io.swagger.client.api.DefaultApi;

ApiClient defaultClient = Configuration.getDefaultApiClient();


DefaultApi apiInstance = new DefaultApi();
RemoveSubjectRequest body = new RemoveSubjectRequest(); // RemoveSubjectRequest | 
try {
    apiInstance.subjectsremovePost(body);
} catch (ApiException e) {
    System.err.println("Exception when calling DefaultApi#subjectsremovePost");
    e.printStackTrace();
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RemoveSubjectRequest**](RemoveSubjectRequest.md)|  |

### Return type

null (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

