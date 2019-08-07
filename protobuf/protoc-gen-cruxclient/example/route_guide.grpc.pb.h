// Generated by the gRPC C++ plugin.
// If you make any local change, they will be lost.
// source: route_guide.proto
// Original file comments:
// Copyright 2015 gRPC authors.
// Copied from https://github.com/grpc/grpc/blob/v1.20.0/examples/protos/route_guide.proto
//
#ifndef GRPC_route_5fguide_2eproto__INCLUDED
#define GRPC_route_5fguide_2eproto__INCLUDED

#include "route_guide.pb.h"

#include <grpcpp/impl/codegen/async_stream.h>
#include <grpcpp/impl/codegen/async_unary_call.h>
#include <grpcpp/impl/codegen/method_handler_impl.h>
#include <grpcpp/impl/codegen/proto_utils.h>
#include <grpcpp/impl/codegen/rpc_method.h>
#include <grpcpp/impl/codegen/service_type.h>
#include <grpcpp/impl/codegen/status.h>
#include <grpcpp/impl/codegen/stub_options.h>
#include <grpcpp/impl/codegen/sync_stream.h>

namespace grpc {
class CompletionQueue;
class Channel;
class ServerCompletionQueue;
class ServerContext;
}  // namespace grpc

namespace routeguide {

// Interface exported by the server.
class RouteGuide final {
 public:
  static constexpr char const* service_full_name() {
    return "routeguide.RouteGuide";
  }
  class StubInterface {
   public:
    virtual ~StubInterface() {}
    // A simple RPC.
    //
    // Obtains the feature at a given position.
    //
    // A feature with an empty name is returned if there's no feature at the given
    // position.
    virtual ::grpc::Status GetFeature(::grpc::ClientContext* context, const ::routeguide::Point& request, ::routeguide::Feature* response) = 0;
    std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::routeguide::Feature>> AsyncGetFeature(::grpc::ClientContext* context, const ::routeguide::Point& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::routeguide::Feature>>(AsyncGetFeatureRaw(context, request, cq));
    }
    std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::routeguide::Feature>> PrepareAsyncGetFeature(::grpc::ClientContext* context, const ::routeguide::Point& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::routeguide::Feature>>(PrepareAsyncGetFeatureRaw(context, request, cq));
    }
    // A server-to-client streaming RPC.
    //
    // Obtains the Features available within the given Rectangle.  Results are
    // streamed rather than returned at once (e.g. in a response message with a
    // repeated field), as the rectangle may cover a large area and contain a
    // huge number of features.
    std::unique_ptr< ::grpc::ClientReaderInterface< ::routeguide::Feature>> ListFeatures(::grpc::ClientContext* context, const ::routeguide::Rectangle& request) {
      return std::unique_ptr< ::grpc::ClientReaderInterface< ::routeguide::Feature>>(ListFeaturesRaw(context, request));
    }
    std::unique_ptr< ::grpc::ClientAsyncReaderInterface< ::routeguide::Feature>> AsyncListFeatures(::grpc::ClientContext* context, const ::routeguide::Rectangle& request, ::grpc::CompletionQueue* cq, void* tag) {
      return std::unique_ptr< ::grpc::ClientAsyncReaderInterface< ::routeguide::Feature>>(AsyncListFeaturesRaw(context, request, cq, tag));
    }
    std::unique_ptr< ::grpc::ClientAsyncReaderInterface< ::routeguide::Feature>> PrepareAsyncListFeatures(::grpc::ClientContext* context, const ::routeguide::Rectangle& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncReaderInterface< ::routeguide::Feature>>(PrepareAsyncListFeaturesRaw(context, request, cq));
    }
    // A client-to-server streaming RPC.
    //
    // Accepts a stream of Points on a route being traversed, returning a
    // RouteSummary when traversal is completed.
    std::unique_ptr< ::grpc::ClientWriterInterface< ::routeguide::Point>> RecordRoute(::grpc::ClientContext* context, ::routeguide::RouteSummary* response) {
      return std::unique_ptr< ::grpc::ClientWriterInterface< ::routeguide::Point>>(RecordRouteRaw(context, response));
    }
    std::unique_ptr< ::grpc::ClientAsyncWriterInterface< ::routeguide::Point>> AsyncRecordRoute(::grpc::ClientContext* context, ::routeguide::RouteSummary* response, ::grpc::CompletionQueue* cq, void* tag) {
      return std::unique_ptr< ::grpc::ClientAsyncWriterInterface< ::routeguide::Point>>(AsyncRecordRouteRaw(context, response, cq, tag));
    }
    std::unique_ptr< ::grpc::ClientAsyncWriterInterface< ::routeguide::Point>> PrepareAsyncRecordRoute(::grpc::ClientContext* context, ::routeguide::RouteSummary* response, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncWriterInterface< ::routeguide::Point>>(PrepareAsyncRecordRouteRaw(context, response, cq));
    }
    // A Bidirectional streaming RPC.
    //
    // Accepts a stream of RouteNotes sent while a route is being traversed,
    // while receiving other RouteNotes (e.g. from other users).
    std::unique_ptr< ::grpc::ClientReaderWriterInterface< ::routeguide::RouteNote, ::routeguide::RouteNote>> RouteChat(::grpc::ClientContext* context) {
      return std::unique_ptr< ::grpc::ClientReaderWriterInterface< ::routeguide::RouteNote, ::routeguide::RouteNote>>(RouteChatRaw(context));
    }
    std::unique_ptr< ::grpc::ClientAsyncReaderWriterInterface< ::routeguide::RouteNote, ::routeguide::RouteNote>> AsyncRouteChat(::grpc::ClientContext* context, ::grpc::CompletionQueue* cq, void* tag) {
      return std::unique_ptr< ::grpc::ClientAsyncReaderWriterInterface< ::routeguide::RouteNote, ::routeguide::RouteNote>>(AsyncRouteChatRaw(context, cq, tag));
    }
    std::unique_ptr< ::grpc::ClientAsyncReaderWriterInterface< ::routeguide::RouteNote, ::routeguide::RouteNote>> PrepareAsyncRouteChat(::grpc::ClientContext* context, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncReaderWriterInterface< ::routeguide::RouteNote, ::routeguide::RouteNote>>(PrepareAsyncRouteChatRaw(context, cq));
    }
  private:
    virtual ::grpc::ClientAsyncResponseReaderInterface< ::routeguide::Feature>* AsyncGetFeatureRaw(::grpc::ClientContext* context, const ::routeguide::Point& request, ::grpc::CompletionQueue* cq) = 0;
    virtual ::grpc::ClientAsyncResponseReaderInterface< ::routeguide::Feature>* PrepareAsyncGetFeatureRaw(::grpc::ClientContext* context, const ::routeguide::Point& request, ::grpc::CompletionQueue* cq) = 0;
    virtual ::grpc::ClientReaderInterface< ::routeguide::Feature>* ListFeaturesRaw(::grpc::ClientContext* context, const ::routeguide::Rectangle& request) = 0;
    virtual ::grpc::ClientAsyncReaderInterface< ::routeguide::Feature>* AsyncListFeaturesRaw(::grpc::ClientContext* context, const ::routeguide::Rectangle& request, ::grpc::CompletionQueue* cq, void* tag) = 0;
    virtual ::grpc::ClientAsyncReaderInterface< ::routeguide::Feature>* PrepareAsyncListFeaturesRaw(::grpc::ClientContext* context, const ::routeguide::Rectangle& request, ::grpc::CompletionQueue* cq) = 0;
    virtual ::grpc::ClientWriterInterface< ::routeguide::Point>* RecordRouteRaw(::grpc::ClientContext* context, ::routeguide::RouteSummary* response) = 0;
    virtual ::grpc::ClientAsyncWriterInterface< ::routeguide::Point>* AsyncRecordRouteRaw(::grpc::ClientContext* context, ::routeguide::RouteSummary* response, ::grpc::CompletionQueue* cq, void* tag) = 0;
    virtual ::grpc::ClientAsyncWriterInterface< ::routeguide::Point>* PrepareAsyncRecordRouteRaw(::grpc::ClientContext* context, ::routeguide::RouteSummary* response, ::grpc::CompletionQueue* cq) = 0;
    virtual ::grpc::ClientReaderWriterInterface< ::routeguide::RouteNote, ::routeguide::RouteNote>* RouteChatRaw(::grpc::ClientContext* context) = 0;
    virtual ::grpc::ClientAsyncReaderWriterInterface< ::routeguide::RouteNote, ::routeguide::RouteNote>* AsyncRouteChatRaw(::grpc::ClientContext* context, ::grpc::CompletionQueue* cq, void* tag) = 0;
    virtual ::grpc::ClientAsyncReaderWriterInterface< ::routeguide::RouteNote, ::routeguide::RouteNote>* PrepareAsyncRouteChatRaw(::grpc::ClientContext* context, ::grpc::CompletionQueue* cq) = 0;
  };
  class Stub final : public StubInterface {
   public:
    Stub(const std::shared_ptr< ::grpc::ChannelInterface>& channel);
    ::grpc::Status GetFeature(::grpc::ClientContext* context, const ::routeguide::Point& request, ::routeguide::Feature* response) override;
    std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::routeguide::Feature>> AsyncGetFeature(::grpc::ClientContext* context, const ::routeguide::Point& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::routeguide::Feature>>(AsyncGetFeatureRaw(context, request, cq));
    }
    std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::routeguide::Feature>> PrepareAsyncGetFeature(::grpc::ClientContext* context, const ::routeguide::Point& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::routeguide::Feature>>(PrepareAsyncGetFeatureRaw(context, request, cq));
    }
    std::unique_ptr< ::grpc::ClientReader< ::routeguide::Feature>> ListFeatures(::grpc::ClientContext* context, const ::routeguide::Rectangle& request) {
      return std::unique_ptr< ::grpc::ClientReader< ::routeguide::Feature>>(ListFeaturesRaw(context, request));
    }
    std::unique_ptr< ::grpc::ClientAsyncReader< ::routeguide::Feature>> AsyncListFeatures(::grpc::ClientContext* context, const ::routeguide::Rectangle& request, ::grpc::CompletionQueue* cq, void* tag) {
      return std::unique_ptr< ::grpc::ClientAsyncReader< ::routeguide::Feature>>(AsyncListFeaturesRaw(context, request, cq, tag));
    }
    std::unique_ptr< ::grpc::ClientAsyncReader< ::routeguide::Feature>> PrepareAsyncListFeatures(::grpc::ClientContext* context, const ::routeguide::Rectangle& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncReader< ::routeguide::Feature>>(PrepareAsyncListFeaturesRaw(context, request, cq));
    }
    std::unique_ptr< ::grpc::ClientWriter< ::routeguide::Point>> RecordRoute(::grpc::ClientContext* context, ::routeguide::RouteSummary* response) {
      return std::unique_ptr< ::grpc::ClientWriter< ::routeguide::Point>>(RecordRouteRaw(context, response));
    }
    std::unique_ptr< ::grpc::ClientAsyncWriter< ::routeguide::Point>> AsyncRecordRoute(::grpc::ClientContext* context, ::routeguide::RouteSummary* response, ::grpc::CompletionQueue* cq, void* tag) {
      return std::unique_ptr< ::grpc::ClientAsyncWriter< ::routeguide::Point>>(AsyncRecordRouteRaw(context, response, cq, tag));
    }
    std::unique_ptr< ::grpc::ClientAsyncWriter< ::routeguide::Point>> PrepareAsyncRecordRoute(::grpc::ClientContext* context, ::routeguide::RouteSummary* response, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncWriter< ::routeguide::Point>>(PrepareAsyncRecordRouteRaw(context, response, cq));
    }
    std::unique_ptr< ::grpc::ClientReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>> RouteChat(::grpc::ClientContext* context) {
      return std::unique_ptr< ::grpc::ClientReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>>(RouteChatRaw(context));
    }
    std::unique_ptr<  ::grpc::ClientAsyncReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>> AsyncRouteChat(::grpc::ClientContext* context, ::grpc::CompletionQueue* cq, void* tag) {
      return std::unique_ptr< ::grpc::ClientAsyncReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>>(AsyncRouteChatRaw(context, cq, tag));
    }
    std::unique_ptr<  ::grpc::ClientAsyncReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>> PrepareAsyncRouteChat(::grpc::ClientContext* context, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>>(PrepareAsyncRouteChatRaw(context, cq));
    }

   private:
    std::shared_ptr< ::grpc::ChannelInterface> channel_;
    ::grpc::ClientAsyncResponseReader< ::routeguide::Feature>* AsyncGetFeatureRaw(::grpc::ClientContext* context, const ::routeguide::Point& request, ::grpc::CompletionQueue* cq) override;
    ::grpc::ClientAsyncResponseReader< ::routeguide::Feature>* PrepareAsyncGetFeatureRaw(::grpc::ClientContext* context, const ::routeguide::Point& request, ::grpc::CompletionQueue* cq) override;
    ::grpc::ClientReader< ::routeguide::Feature>* ListFeaturesRaw(::grpc::ClientContext* context, const ::routeguide::Rectangle& request) override;
    ::grpc::ClientAsyncReader< ::routeguide::Feature>* AsyncListFeaturesRaw(::grpc::ClientContext* context, const ::routeguide::Rectangle& request, ::grpc::CompletionQueue* cq, void* tag) override;
    ::grpc::ClientAsyncReader< ::routeguide::Feature>* PrepareAsyncListFeaturesRaw(::grpc::ClientContext* context, const ::routeguide::Rectangle& request, ::grpc::CompletionQueue* cq) override;
    ::grpc::ClientWriter< ::routeguide::Point>* RecordRouteRaw(::grpc::ClientContext* context, ::routeguide::RouteSummary* response) override;
    ::grpc::ClientAsyncWriter< ::routeguide::Point>* AsyncRecordRouteRaw(::grpc::ClientContext* context, ::routeguide::RouteSummary* response, ::grpc::CompletionQueue* cq, void* tag) override;
    ::grpc::ClientAsyncWriter< ::routeguide::Point>* PrepareAsyncRecordRouteRaw(::grpc::ClientContext* context, ::routeguide::RouteSummary* response, ::grpc::CompletionQueue* cq) override;
    ::grpc::ClientReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>* RouteChatRaw(::grpc::ClientContext* context) override;
    ::grpc::ClientAsyncReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>* AsyncRouteChatRaw(::grpc::ClientContext* context, ::grpc::CompletionQueue* cq, void* tag) override;
    ::grpc::ClientAsyncReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>* PrepareAsyncRouteChatRaw(::grpc::ClientContext* context, ::grpc::CompletionQueue* cq) override;
    const ::grpc::internal::RpcMethod rpcmethod_GetFeature_;
    const ::grpc::internal::RpcMethod rpcmethod_ListFeatures_;
    const ::grpc::internal::RpcMethod rpcmethod_RecordRoute_;
    const ::grpc::internal::RpcMethod rpcmethod_RouteChat_;
  };
  static std::unique_ptr<Stub> NewStub(const std::shared_ptr< ::grpc::ChannelInterface>& channel, const ::grpc::StubOptions& options = ::grpc::StubOptions());

  class Service : public ::grpc::Service {
   public:
    Service();
    virtual ~Service();
    // A simple RPC.
    //
    // Obtains the feature at a given position.
    //
    // A feature with an empty name is returned if there's no feature at the given
    // position.
    virtual ::grpc::Status GetFeature(::grpc::ServerContext* context, const ::routeguide::Point* request, ::routeguide::Feature* response);
    // A server-to-client streaming RPC.
    //
    // Obtains the Features available within the given Rectangle.  Results are
    // streamed rather than returned at once (e.g. in a response message with a
    // repeated field), as the rectangle may cover a large area and contain a
    // huge number of features.
    virtual ::grpc::Status ListFeatures(::grpc::ServerContext* context, const ::routeguide::Rectangle* request, ::grpc::ServerWriter< ::routeguide::Feature>* writer);
    // A client-to-server streaming RPC.
    //
    // Accepts a stream of Points on a route being traversed, returning a
    // RouteSummary when traversal is completed.
    virtual ::grpc::Status RecordRoute(::grpc::ServerContext* context, ::grpc::ServerReader< ::routeguide::Point>* reader, ::routeguide::RouteSummary* response);
    // A Bidirectional streaming RPC.
    //
    // Accepts a stream of RouteNotes sent while a route is being traversed,
    // while receiving other RouteNotes (e.g. from other users).
    virtual ::grpc::Status RouteChat(::grpc::ServerContext* context, ::grpc::ServerReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>* stream);
  };
  template <class BaseClass>
  class WithAsyncMethod_GetFeature : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithAsyncMethod_GetFeature() {
      ::grpc::Service::MarkMethodAsync(0);
    }
    ~WithAsyncMethod_GetFeature() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status GetFeature(::grpc::ServerContext* context, const ::routeguide::Point* request, ::routeguide::Feature* response) final override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    void RequestGetFeature(::grpc::ServerContext* context, ::routeguide::Point* request, ::grpc::ServerAsyncResponseWriter< ::routeguide::Feature>* response, ::grpc::CompletionQueue* new_call_cq, ::grpc::ServerCompletionQueue* notification_cq, void *tag) {
      ::grpc::Service::RequestAsyncUnary(0, context, request, response, new_call_cq, notification_cq, tag);
    }
  };
  template <class BaseClass>
  class WithAsyncMethod_ListFeatures : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithAsyncMethod_ListFeatures() {
      ::grpc::Service::MarkMethodAsync(1);
    }
    ~WithAsyncMethod_ListFeatures() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status ListFeatures(::grpc::ServerContext* context, const ::routeguide::Rectangle* request, ::grpc::ServerWriter< ::routeguide::Feature>* writer) final override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    void RequestListFeatures(::grpc::ServerContext* context, ::routeguide::Rectangle* request, ::grpc::ServerAsyncWriter< ::routeguide::Feature>* writer, ::grpc::CompletionQueue* new_call_cq, ::grpc::ServerCompletionQueue* notification_cq, void *tag) {
      ::grpc::Service::RequestAsyncServerStreaming(1, context, request, writer, new_call_cq, notification_cq, tag);
    }
  };
  template <class BaseClass>
  class WithAsyncMethod_RecordRoute : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithAsyncMethod_RecordRoute() {
      ::grpc::Service::MarkMethodAsync(2);
    }
    ~WithAsyncMethod_RecordRoute() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status RecordRoute(::grpc::ServerContext* context, ::grpc::ServerReader< ::routeguide::Point>* reader, ::routeguide::RouteSummary* response) final override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    void RequestRecordRoute(::grpc::ServerContext* context, ::grpc::ServerAsyncReader< ::routeguide::RouteSummary, ::routeguide::Point>* reader, ::grpc::CompletionQueue* new_call_cq, ::grpc::ServerCompletionQueue* notification_cq, void *tag) {
      ::grpc::Service::RequestAsyncClientStreaming(2, context, reader, new_call_cq, notification_cq, tag);
    }
  };
  template <class BaseClass>
  class WithAsyncMethod_RouteChat : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithAsyncMethod_RouteChat() {
      ::grpc::Service::MarkMethodAsync(3);
    }
    ~WithAsyncMethod_RouteChat() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status RouteChat(::grpc::ServerContext* context, ::grpc::ServerReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>* stream) final override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    void RequestRouteChat(::grpc::ServerContext* context, ::grpc::ServerAsyncReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>* stream, ::grpc::CompletionQueue* new_call_cq, ::grpc::ServerCompletionQueue* notification_cq, void *tag) {
      ::grpc::Service::RequestAsyncBidiStreaming(3, context, stream, new_call_cq, notification_cq, tag);
    }
  };
  typedef WithAsyncMethod_GetFeature<WithAsyncMethod_ListFeatures<WithAsyncMethod_RecordRoute<WithAsyncMethod_RouteChat<Service > > > > AsyncService;
  template <class BaseClass>
  class WithGenericMethod_GetFeature : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithGenericMethod_GetFeature() {
      ::grpc::Service::MarkMethodGeneric(0);
    }
    ~WithGenericMethod_GetFeature() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status GetFeature(::grpc::ServerContext* context, const ::routeguide::Point* request, ::routeguide::Feature* response) final override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
  };
  template <class BaseClass>
  class WithGenericMethod_ListFeatures : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithGenericMethod_ListFeatures() {
      ::grpc::Service::MarkMethodGeneric(1);
    }
    ~WithGenericMethod_ListFeatures() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status ListFeatures(::grpc::ServerContext* context, const ::routeguide::Rectangle* request, ::grpc::ServerWriter< ::routeguide::Feature>* writer) final override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
  };
  template <class BaseClass>
  class WithGenericMethod_RecordRoute : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithGenericMethod_RecordRoute() {
      ::grpc::Service::MarkMethodGeneric(2);
    }
    ~WithGenericMethod_RecordRoute() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status RecordRoute(::grpc::ServerContext* context, ::grpc::ServerReader< ::routeguide::Point>* reader, ::routeguide::RouteSummary* response) final override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
  };
  template <class BaseClass>
  class WithGenericMethod_RouteChat : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithGenericMethod_RouteChat() {
      ::grpc::Service::MarkMethodGeneric(3);
    }
    ~WithGenericMethod_RouteChat() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status RouteChat(::grpc::ServerContext* context, ::grpc::ServerReaderWriter< ::routeguide::RouteNote, ::routeguide::RouteNote>* stream) final override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
  };
  template <class BaseClass>
  class WithStreamedUnaryMethod_GetFeature : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithStreamedUnaryMethod_GetFeature() {
      ::grpc::Service::MarkMethodStreamed(0,
        new ::grpc::internal::StreamedUnaryHandler< ::routeguide::Point, ::routeguide::Feature>(std::bind(&WithStreamedUnaryMethod_GetFeature<BaseClass>::StreamedGetFeature, this, std::placeholders::_1, std::placeholders::_2)));
    }
    ~WithStreamedUnaryMethod_GetFeature() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable regular version of this method
    ::grpc::Status GetFeature(::grpc::ServerContext* context, const ::routeguide::Point* request, ::routeguide::Feature* response) final override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    // replace default version of method with streamed unary
    virtual ::grpc::Status StreamedGetFeature(::grpc::ServerContext* context, ::grpc::ServerUnaryStreamer< ::routeguide::Point,::routeguide::Feature>* server_unary_streamer) = 0;
  };
  typedef WithStreamedUnaryMethod_GetFeature<Service > StreamedUnaryService;
  template <class BaseClass>
  class WithSplitStreamingMethod_ListFeatures : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithSplitStreamingMethod_ListFeatures() {
      ::grpc::Service::MarkMethodStreamed(1,
        new ::grpc::internal::SplitServerStreamingHandler< ::routeguide::Rectangle, ::routeguide::Feature>(std::bind(&WithSplitStreamingMethod_ListFeatures<BaseClass>::StreamedListFeatures, this, std::placeholders::_1, std::placeholders::_2)));
    }
    ~WithSplitStreamingMethod_ListFeatures() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable regular version of this method
    ::grpc::Status ListFeatures(::grpc::ServerContext* context, const ::routeguide::Rectangle* request, ::grpc::ServerWriter< ::routeguide::Feature>* writer) final override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    // replace default version of method with split streamed
    virtual ::grpc::Status StreamedListFeatures(::grpc::ServerContext* context, ::grpc::ServerSplitStreamer< ::routeguide::Rectangle,::routeguide::Feature>* server_split_streamer) = 0;
  };
  typedef WithSplitStreamingMethod_ListFeatures<Service > SplitStreamedService;
  typedef WithStreamedUnaryMethod_GetFeature<WithSplitStreamingMethod_ListFeatures<Service > > StreamedService;
};

}  // namespace routeguide


#endif  // GRPC_route_5fguide_2eproto__INCLUDED
