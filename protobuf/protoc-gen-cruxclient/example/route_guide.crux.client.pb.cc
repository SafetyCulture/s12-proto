// Generated by the CRUX Service Client C++ plugin.
// If you make any local change, they will be lost.
// source: route_guide.proto

#include "route_guide.crux.client.pb.h"
#include "s12_client_support.hpp"

namespace routeguide {

RouteGuideClient::RouteGuideClient(const std::shared_ptr<RouteGuide::StubInterface>& stub) : mStub(stub) {}

void RouteGuideClient::Invoke(const google::protobuf::Any& request_data, const std::string& method) const {
  if (method == kRouteGuideGetFeature) {
    routeguide::Point request;
    if (!request_data.UnpackTo(&request)) {
      throw crux::RequestParseException();
    }
    GetFeature(request);
    return;
  } else if (method == kRouteGuideUpdateFeature) {
    routeguide::Point request;
    if (!request_data.UnpackTo(&request)) {
      throw crux::RequestParseException();
    }
    UpdateFeature(request);
    return;
  } else if (method == kRouteGuideListFeatures) {
    routeguide::Rectangle request;
    if (!request_data.UnpackTo(&request)) {
      throw crux::RequestParseException();
    }
    ListFeatures(request);
    return;
  }
  throw crux::RequestParseException();
}
routeguide::Feature RouteGuideClient::GetFeature(const routeguide::Point& request) const {
  routeguide::Feature response;
  grpc::ClientContext context;
  grpc::Status status = mStub->GetFeature(&context, request, &response);
  if (!status.ok()) {
    throw crux::ServiceException(status.error_code(), status.error_message());
  }
  return response;
}
routeguide::Feature RouteGuideClient::UpdateFeature(const routeguide::Point& request) const {
  routeguide::Feature response;
  grpc::ClientContext context;
  grpc::Status status = mStub->UpdateFeature(&context, request, &response);
  if (!status.ok()) {
    throw crux::ServiceException(status.error_code(), status.error_message());
  }
  return response;
}
std::vector<routeguide::Feature> RouteGuideClient::ListFeatures(const routeguide::Rectangle& request) const {
  std::vector<routeguide::Feature> response;
  grpc::ClientContext context;
  routeguide::Feature item;
  std::unique_ptr<grpc::ClientReaderInterface<routeguide::Feature>> stream = mStub->ListFeatures(&context, request);
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
    routeguide::Point request;
    if (!request_data.UnpackTo(&request)) {
      throw crux::RequestParseException();
    }
    GetFeature(request);
    return;
  }
  throw crux::RequestParseException();
}
routeguide::Feature PublicRouteGuideClient::GetFeature(const routeguide::Point& request) const {
  routeguide::Feature response;
  grpc::ClientContext context;
  grpc::Status status = mStub->GetFeature(&context, request, &response);
  if (!status.ok()) {
    throw crux::ServiceException(status.error_code(), status.error_message());
  }
  return response;
}

}

