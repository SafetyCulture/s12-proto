// Generated by the CRUX Service Client C++ plugin.
// If you make any local change, they will be lost.
// source: route_guide.proto
#pragma once

#include <string>
#include <memory>

#include <google/protobuf/any.pb.h>
#include "route_guide.grpc.pb.h"

namespace routeguide {

class RouteGuideClientInterface {
 public:
  virtual ~RouteGuideClientInterface() {}
  virtual void Invoke(const google::protobuf::Any& request_data) const = 0;
  virtual routeguide::Feature GetFeature(const routeguide::Point& request) const = 0;
  virtual std::vector<routeguide::Feature> ListFeatures(const routeguide::Rectangle& request) const = 0;
};

class RouteGuideClient: public RouteGuideClientInterface {
 public:
  explicit RouteGuideClient(const std::shared_ptr<RouteGuide::StubInterface>& stub);
  void Invoke(const google::protobuf::Any& request_data) const override;
  routeguide::Feature GetFeature(const routeguide::Point& request) const override;
  std::vector<routeguide::Feature> ListFeatures(const routeguide::Rectangle& request) const override;

 private:
  std::shared_ptr<RouteGuide::StubInterface> mStub;

};

}

