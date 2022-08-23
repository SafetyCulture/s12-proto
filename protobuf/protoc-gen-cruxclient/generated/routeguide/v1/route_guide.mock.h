// Generated by the CRUX Engine C++ plugin.
// If you make any local change, they will be lost.
// source: routeguide/v1/route_guide.proto
#pragma once

#include <tuple>
#include "routeguide/v1/route_guide.grpc.pb.h"

#include "crux_mock_support.h"

namespace routeguide::v1 {
class RouteGuideImpl final : public ::routeguide::v1::RouteGuide::Service {
 public:
  RouteGuideImpl(const GRPCServiceCallback& callback): mCallback(callback) {};
  ::grpc::Status GetFeature(::grpc::ServerContext* context, const ::routeguide::v1::Point* request, ::routeguide::v1::Feature* response) override {
    auto[status, bytes] = mCallback("routeguide.v1.RouteGuide.GetFeature", request->SerializeAsString(), context->client_metadata());
    if (bytes.has_value()) {
      response->ParseFromString(*bytes);
    }
    return status;
  }
  ::grpc::Status UpdateFeature(::grpc::ServerContext* context, const ::routeguide::v1::Point* request, ::routeguide::v1::Feature* response) override {
    auto[status, bytes] = mCallback("routeguide.v1.RouteGuide.UpdateFeature", request->SerializeAsString(), context->client_metadata());
    if (bytes.has_value()) {
      response->ParseFromString(*bytes);
    }
    return status;
  }
  ::grpc::Status ListFeatures(::grpc::ServerContext* context, const ::routeguide::v1::Rectangle* request, ::routeguide::v1::Feature* response) override {
    auto[status, bytes] = mCallback("routeguide.v1.RouteGuide.ListFeatures", request->SerializeAsString(), context->client_metadata());
    if (bytes.has_value()) {
      response->ParseFromString(*bytes);
    }
    return status;
  }
 private:
  const GRPCServiceCallback mCallback;
};

class PublicRouteGuideImpl final : public ::routeguide::v1::PublicRouteGuide::Service {
 public:
  PublicRouteGuideImpl(const GRPCServiceCallback& callback): mCallback(callback) {};
  ::grpc::Status GetFeature(::grpc::ServerContext* context, const ::routeguide::v1::Point* request, ::routeguide::v1::Feature* response) override {
    auto[status, bytes] = mCallback("routeguide.v1.PublicRouteGuide.GetFeature", request->SerializeAsString(), context->client_metadata());
    if (bytes.has_value()) {
      response->ParseFromString(*bytes);
    }
    return status;
  }
 private:
  const GRPCServiceCallback mCallback;
};

}  // namespace routeguide::v1
