// Generated by the CRUX Engine C++ plugin.
// If you make any local change, they will be lost.
// source: routeguide/v1/route_guide.proto
#pragma once

#include <tuple>
#include "routeguide/v1/route_guide.grpc.pb.h"

#include "crux_mock_support.h"

namespace routeguide::v1 {
class RouteGuideMockImpl final : public ::routeguide::v1::RouteGuide::Service {
 public:
  RouteGuideMockImpl(const GRPCServiceCallback& callback): mCallback(callback) {};
  ::grpc::Status GetFeature(::grpc::ServerContext* context, const ::routeguide::v1::Point* request, ::routeguide::v1::Feature* response) override {
    auto[status, bytes_list] = mCallback("routeguide.v1.RouteGuide.GetFeature", request->SerializeAsString(), context->client_metadata());
    if (!bytes_list.empty()) {
      response->ParseFromString(bytes_list[0]);
    }
    return status;
  }
  ::grpc::Status UpdateFeature(::grpc::ServerContext* context, const ::routeguide::v1::Point* request, ::routeguide::v1::Feature* response) override {
    auto[status, bytes_list] = mCallback("routeguide.v1.RouteGuide.UpdateFeature", request->SerializeAsString(), context->client_metadata());
    if (!bytes_list.empty()) {
      response->ParseFromString(bytes_list[0]);
    }
    return status;
  }
  ::grpc::Status ListFeatures(::grpc::ServerContext* context, const ::routeguide::v1::Rectangle* request, ::grpc::ServerWriter<::routeguide::v1::Feature>* writer) override {
    auto[status, bytes_list] = mCallback("routeguide.v1.RouteGuide.ListFeatures", request->SerializeAsString(), context->client_metadata());
    for (auto& bytes : bytes_list) {
      ::routeguide::v1::Feature response;
      response.ParseFromString(bytes);
      writer->Write(response);
    }
    return status;
  }
 private:
  const GRPCServiceCallback mCallback;
};

class PublicRouteGuideMockImpl final : public ::routeguide::v1::PublicRouteGuide::Service {
 public:
  PublicRouteGuideMockImpl(const GRPCServiceCallback& callback): mCallback(callback) {};
  ::grpc::Status GetFeature(::grpc::ServerContext* context, const ::routeguide::v1::Point* request, ::routeguide::v1::Feature* response) override {
    auto[status, bytes_list] = mCallback("routeguide.v1.PublicRouteGuide.GetFeature", request->SerializeAsString(), context->client_metadata());
    if (!bytes_list.empty()) {
      response->ParseFromString(bytes_list[0]);
    }
    return status;
  }
 private:
  const GRPCServiceCallback mCallback;
};

}  // namespace routeguide::v1
