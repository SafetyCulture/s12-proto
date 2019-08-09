// Generated by the CRUX Service Client C++ plugin.
// If you make any local change, they will be lost.
// source: route_guide.proto
#pragma once

#include <vector>
#include <string>
#include <memory>

#include <google/protobuf/any.pb.h>
#include "route_guide.grpc.pb.h"
#include "s12_client_support.hpp"

namespace routeguide {

const char kRouteGuideGetFeature[] = "/routeguide.RouteGuide/GetFeature";
const char kRouteGuideUpdateFeature[] = "/routeguide.RouteGuide/UpdateFeature";
const char kRouteGuideListFeatures[] = "/routeguide.RouteGuide/ListFeatures";
class RouteGuideClientInterface {
 public:
  virtual ~RouteGuideClientInterface() {}
  virtual void Invoke(const google::protobuf::Any& request_data, const std::string& method) const {}
  virtual routeguide::Feature GetFeature(const routeguide::Point& request) const = 0;
  virtual routeguide::Feature UpdateFeature(const routeguide::Point& request) const = 0;
  virtual std::vector<routeguide::Feature> ListFeatures(const routeguide::Rectangle& request) const = 0;
};

const char kPublicRouteGuideGetFeature[] = "/routeguide.PublicRouteGuide/GetFeature";
class PublicRouteGuideClientInterface {
 public:
  virtual ~PublicRouteGuideClientInterface() {}
  virtual void Invoke(const google::protobuf::Any& request_data, const std::string& method) const {}
  virtual routeguide::Feature GetFeature(const routeguide::Point& request) const = 0;
};

class RouteGuideClient: public RouteGuideClientInterface {
 public:
  explicit RouteGuideClient(const std::shared_ptr<RouteGuide::StubInterface>& stub);
  void Invoke(const google::protobuf::Any& request_data, const std::string& method) const override;
  routeguide::Feature GetFeature(const routeguide::Point& request) const override;
  routeguide::Feature UpdateFeature(const routeguide::Point& request) const override;
  std::vector<routeguide::Feature> ListFeatures(const routeguide::Rectangle& request) const override;

 private:
  std::shared_ptr<RouteGuide::StubInterface> mStub;

};

class PublicRouteGuideClient: public PublicRouteGuideClientInterface {
 public:
  explicit PublicRouteGuideClient(const std::shared_ptr<PublicRouteGuide::StubInterface>& stub);
  void Invoke(const google::protobuf::Any& request_data, const std::string& method) const override;
  routeguide::Feature GetFeature(const routeguide::Point& request) const override;

 private:
  std::shared_ptr<PublicRouteGuide::StubInterface> mStub;

};

class MockRouteGuideClient: public RouteGuideClientInterface {
 public:
  mutable int mInvokeCalledCount = 0;
  mutable std::vector<google::protobuf::Any> mInvokeRequestData;
  mutable std::vector<std::string> mInvokeMethods;
  bool mInvokeThrowParseException = false;
  grpc::StatusCode mInvokeErrorStatusCode = grpc::StatusCode::INVALID_ARGUMENT;
  mutable int mInvokeExceptionThrowCount = 0;
  void Invoke(const google::protobuf::Any& request_data, const std::string& method) const override {
    mInvokeCalledCount++;
    mInvokeRequestData.push_back(request_data);
    mInvokeMethods.push_back(method);
    if (mInvokeThrowParseException) {
      throw crux::RequestParseException();
    }
    if (mInvokeExceptionThrowCount > 0) {
      mInvokeExceptionThrowCount--;
      throw crux::ServiceException(mInvokeErrorStatusCode, "Error");
    }
  }
  mutable int mGetFeatureCalledCount = 0;
  mutable std::vector<routeguide::Point> mGetFeatureRequests;
  routeguide::Feature mGetFeatureResponse;
  grpc::StatusCode mGetFeatureErrorStatusCode = grpc::StatusCode::INVALID_ARGUMENT;
  mutable int mGetFeatureExceptionThrowCount = 0;
  routeguide::Feature GetFeature(const routeguide::Point& request) const override {
    mGetFeatureCalledCount++;
    mGetFeatureRequests.push_back(request);
    if (mGetFeatureExceptionThrowCount > 0) {
      mGetFeatureExceptionThrowCount--;
      throw crux::ServiceException(mGetFeatureErrorStatusCode, "Error");
    }
    return mGetFeatureResponse;
  }

  mutable int mUpdateFeatureCalledCount = 0;
  mutable std::vector<routeguide::Point> mUpdateFeatureRequests;
  routeguide::Feature mUpdateFeatureResponse;
  grpc::StatusCode mUpdateFeatureErrorStatusCode = grpc::StatusCode::INVALID_ARGUMENT;
  mutable int mUpdateFeatureExceptionThrowCount = 0;
  routeguide::Feature UpdateFeature(const routeguide::Point& request) const override {
    mUpdateFeatureCalledCount++;
    mUpdateFeatureRequests.push_back(request);
    if (mUpdateFeatureExceptionThrowCount > 0) {
      mUpdateFeatureExceptionThrowCount--;
      throw crux::ServiceException(mUpdateFeatureErrorStatusCode, "Error");
    }
    return mUpdateFeatureResponse;
  }

  mutable int mListFeaturesCalledCount = 0;
  mutable std::vector<routeguide::Rectangle> mListFeaturesRequests;
  std::vector<routeguide::Feature> mListFeaturesResponse;
  grpc::StatusCode mListFeaturesErrorStatusCode = grpc::StatusCode::INVALID_ARGUMENT;
  mutable int mListFeaturesExceptionThrowCount = 0;
  std::vector<routeguide::Feature> ListFeatures(const routeguide::Rectangle& request) const override {
    mListFeaturesCalledCount++;
    mListFeaturesRequests.push_back(request);
    if (mListFeaturesExceptionThrowCount > 0) {
      mListFeaturesExceptionThrowCount--;
      throw crux::ServiceException(mListFeaturesErrorStatusCode, "Error");
    }
    return mListFeaturesResponse;
  }

};

class MockPublicRouteGuideClient: public PublicRouteGuideClientInterface {
 public:
  mutable int mInvokeCalledCount = 0;
  mutable std::vector<google::protobuf::Any> mInvokeRequestData;
  mutable std::vector<std::string> mInvokeMethods;
  bool mInvokeThrowParseException = false;
  grpc::StatusCode mInvokeErrorStatusCode = grpc::StatusCode::INVALID_ARGUMENT;
  mutable int mInvokeExceptionThrowCount = 0;
  void Invoke(const google::protobuf::Any& request_data, const std::string& method) const override {
    mInvokeCalledCount++;
    mInvokeRequestData.push_back(request_data);
    mInvokeMethods.push_back(method);
    if (mInvokeThrowParseException) {
      throw crux::RequestParseException();
    }
    if (mInvokeExceptionThrowCount > 0) {
      mInvokeExceptionThrowCount--;
      throw crux::ServiceException(mInvokeErrorStatusCode, "Error");
    }
  }
  mutable int mGetFeatureCalledCount = 0;
  mutable std::vector<routeguide::Point> mGetFeatureRequests;
  routeguide::Feature mGetFeatureResponse;
  grpc::StatusCode mGetFeatureErrorStatusCode = grpc::StatusCode::INVALID_ARGUMENT;
  mutable int mGetFeatureExceptionThrowCount = 0;
  routeguide::Feature GetFeature(const routeguide::Point& request) const override {
    mGetFeatureCalledCount++;
    mGetFeatureRequests.push_back(request);
    if (mGetFeatureExceptionThrowCount > 0) {
      mGetFeatureExceptionThrowCount--;
      throw crux::ServiceException(mGetFeatureErrorStatusCode, "Error");
    }
    return mGetFeatureResponse;
  }

};

}

