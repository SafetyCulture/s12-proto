// Generated by the CRUX Service Client C++ plugin.
// If you make any local change, they will be lost.
// source: route_guide.proto

#include "route_guide.crux.client.v2.pb.h"
namespace routeguide.v1 {
RouteGuideClient::RouteGuideClient(const std::shared_ptr<RouteGuide::StubInterface>& stub) : mStub(stub) {}

void RouteGuideClient::Invoke(const google::protobuf::Any& request_data, const std::string& method) const {
  if (method == kRouteGuideGetFeature) {
    routeguide::v1::Point request;
    if (!request_data.UnpackTo(&request)) {
      throw crux::RequestParseException();
    }
    GetFeature(request);
    return;
  }
  if (method == kRouteGuideUpdateFeature) {
    routeguide::v1::Point request;
    if (!request_data.UnpackTo(&request)) {
      throw crux::RequestParseException();
    }
    UpdateFeature(request);
    return;
  }
  if (method == kRouteGuideListFeatures) {
    routeguide::v1::Rectangle request;
    if (!request_data.UnpackTo(&request)) {
      throw crux::RequestParseException();
    }
    ListFeatures(request);
    return;
  }
  throw crux::RequestParseException();
}
routeguide::v1::Feature RouteGuideClient::GetFeature(const routeguide::v1::Point& request) const {
  routeguide::v1::Feature response;
  grpc::ClientContext context;
  grpc::Status status = mStub->GetFeature(&context, request, &response);
  if (!status.ok()) {
    throw crux::ServiceException(status.error_code(), status.error_message());
  }
  return response;
}
routeguide::v1::Feature RouteGuideClient::UpdateFeature(const routeguide::v1::Point& request) const {
  routeguide::v1::Feature response;
  grpc::ClientContext context;
  grpc::Status status = mStub->UpdateFeature(&context, request, &response);
  if (!status.ok()) {
    throw crux::ServiceException(status.error_code(), status.error_message());
  }
  return response;
}
std::vector<routeguide::v1::Feature> RouteGuideClient::ListFeatures(const routeguide::v1::Rectangle& request) const {
  std::vector<routeguide::v1::Feature> response;
  grpc::ClientContext context;
  routeguide::v1::Feature item;
  std::unique_ptr<grpc::ClientReaderInterface<routeguide::v1::Feature>> stream = mStub->ListFeatures(&context, request);
  while (stream->Read(&item)) {
    response.emplace_back(item);
  }
  grpc::Status status = stream->Finish();
  if (!status.ok()) {
    throw crux::ServiceException(status.error_code(), status.error_message());
  }
  return response;
}
PublicRouteGuideClient::PublicRouteGuideClient(const std::shared_ptr<PublicRouteGuide::StubInterface>& stub) : mStub(stub) {}

void PublicRouteGuideClient::Invoke(const google::protobuf::Any& request_data, const std::string& method) const {
  if (method == kPublicRouteGuideGetFeature) {
    routeguide::v1::Point request;
    if (!request_data.UnpackTo(&request)) {
      throw crux::RequestParseException();
    }
    GetFeature(request);
    return;
  }
  throw crux::RequestParseException();
}
routeguide::v1::Feature PublicRouteGuideClient::GetFeature(const routeguide::v1::Point& request) const {
  routeguide::v1::Feature response;
  grpc::ClientContext context;
  grpc::Status status = mStub->GetFeature(&context, request, &response);
  if (!status.ok()) {
    throw crux::ServiceException(status.error_code(), status.error_message());
  }
  return response;
}

}  // namespace routeguide.v1
