// GENERATED CODE -- DO NOT EDIT!

// package: 
// file: service.proto

import * as service_pb from "./service_pb";
import * as grpc from "grpc";

interface IOrcaCoreService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  registerProcessor: grpc.MethodDefinition<service_pb.ProcessorRegistration, service_pb.Status>;
  emitWindow: grpc.MethodDefinition<service_pb.Window, service_pb.WindowEmitStatus>;
  readWindowTypes: grpc.MethodDefinition<service_pb.WindowTypeRead, service_pb.WindowTypes>;
}

export const OrcaCoreService: IOrcaCoreService;

export interface IOrcaCoreServer extends grpc.UntypedServiceImplementation {
  registerProcessor: grpc.handleUnaryCall<service_pb.ProcessorRegistration, service_pb.Status>;
  emitWindow: grpc.handleUnaryCall<service_pb.Window, service_pb.WindowEmitStatus>;
  readWindowTypes: grpc.handleUnaryCall<service_pb.WindowTypeRead, service_pb.WindowTypes>;
}

export class OrcaCoreClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  registerProcessor(argument: service_pb.ProcessorRegistration, callback: grpc.requestCallback<service_pb.Status>): grpc.ClientUnaryCall;
  registerProcessor(argument: service_pb.ProcessorRegistration, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<service_pb.Status>): grpc.ClientUnaryCall;
  registerProcessor(argument: service_pb.ProcessorRegistration, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<service_pb.Status>): grpc.ClientUnaryCall;
  emitWindow(argument: service_pb.Window, callback: grpc.requestCallback<service_pb.WindowEmitStatus>): grpc.ClientUnaryCall;
  emitWindow(argument: service_pb.Window, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<service_pb.WindowEmitStatus>): grpc.ClientUnaryCall;
  emitWindow(argument: service_pb.Window, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<service_pb.WindowEmitStatus>): grpc.ClientUnaryCall;
  readWindowTypes(argument: service_pb.WindowTypeRead, callback: grpc.requestCallback<service_pb.WindowTypes>): grpc.ClientUnaryCall;
  readWindowTypes(argument: service_pb.WindowTypeRead, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<service_pb.WindowTypes>): grpc.ClientUnaryCall;
  readWindowTypes(argument: service_pb.WindowTypeRead, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<service_pb.WindowTypes>): grpc.ClientUnaryCall;
}

interface IOrcaProcessorService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  executeDagPart: grpc.MethodDefinition<service_pb.ExecutionRequest, service_pb.ExecutionResult>;
  healthCheck: grpc.MethodDefinition<service_pb.HealthCheckRequest, service_pb.HealthCheckResponse>;
}

export const OrcaProcessorService: IOrcaProcessorService;

export interface IOrcaProcessorServer extends grpc.UntypedServiceImplementation {
  executeDagPart: grpc.handleServerStreamingCall<service_pb.ExecutionRequest, service_pb.ExecutionResult>;
  healthCheck: grpc.handleUnaryCall<service_pb.HealthCheckRequest, service_pb.HealthCheckResponse>;
}

export class OrcaProcessorClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  executeDagPart(argument: service_pb.ExecutionRequest, metadataOrOptions?: grpc.Metadata | grpc.CallOptions | null): grpc.ClientReadableStream<service_pb.ExecutionResult>;
  executeDagPart(argument: service_pb.ExecutionRequest, metadata?: grpc.Metadata | null, options?: grpc.CallOptions | null): grpc.ClientReadableStream<service_pb.ExecutionResult>;
  healthCheck(argument: service_pb.HealthCheckRequest, callback: grpc.requestCallback<service_pb.HealthCheckResponse>): grpc.ClientUnaryCall;
  healthCheck(argument: service_pb.HealthCheckRequest, metadataOrOptions: grpc.Metadata | grpc.CallOptions | null, callback: grpc.requestCallback<service_pb.HealthCheckResponse>): grpc.ClientUnaryCall;
  healthCheck(argument: service_pb.HealthCheckRequest, metadata: grpc.Metadata | null, options: grpc.CallOptions | null, callback: grpc.requestCallback<service_pb.HealthCheckResponse>): grpc.ClientUnaryCall;
}
