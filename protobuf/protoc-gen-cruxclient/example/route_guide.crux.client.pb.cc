// Generated by the CRUX Service Client C++ plugin.
// If you make any local change, they will be lost.
// source: route_guide.proto

#include "route_guide.crux.client.pb.h"

crux::RouteGuideClient::RouteGuideClient(const std::shared_ptr<routeguide::RouteGuide::StubInterface>& stub) : mStub(stub) {}

routeguide::Feature crux::RouteGuideClient::GetFeature(const routeguide::Point& request) const {
  routeguide::Feature response;
  auto status = MakeRequest([stub = mStub, request, &response](){
    grpc::ClientContext context;
    return stub->GetFeature(&context, request, &response);
  });
  if (!status.ok()) {
    throw ServiceException(status.error_code(), status.error_message());
  }
  return response;
}
