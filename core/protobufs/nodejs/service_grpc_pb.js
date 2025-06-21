// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var service_pb = require('./service_pb.js');
var google_protobuf_struct_pb = require('google-protobuf/google/protobuf/struct_pb.js');
var vendor_validate_pb = require('./vendor/validate_pb.js');

function serialize_ExecutionRequest(arg) {
  if (!(arg instanceof service_pb.ExecutionRequest)) {
    throw new Error('Expected argument of type ExecutionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_ExecutionRequest(buffer_arg) {
  return service_pb.ExecutionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_ExecutionResult(arg) {
  if (!(arg instanceof service_pb.ExecutionResult)) {
    throw new Error('Expected argument of type ExecutionResult');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_ExecutionResult(buffer_arg) {
  return service_pb.ExecutionResult.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_HealthCheckRequest(arg) {
  if (!(arg instanceof service_pb.HealthCheckRequest)) {
    throw new Error('Expected argument of type HealthCheckRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_HealthCheckRequest(buffer_arg) {
  return service_pb.HealthCheckRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_HealthCheckResponse(arg) {
  if (!(arg instanceof service_pb.HealthCheckResponse)) {
    throw new Error('Expected argument of type HealthCheckResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_HealthCheckResponse(buffer_arg) {
  return service_pb.HealthCheckResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_ProcessorRegistration(arg) {
  if (!(arg instanceof service_pb.ProcessorRegistration)) {
    throw new Error('Expected argument of type ProcessorRegistration');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_ProcessorRegistration(buffer_arg) {
  return service_pb.ProcessorRegistration.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_Status(arg) {
  if (!(arg instanceof service_pb.Status)) {
    throw new Error('Expected argument of type Status');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_Status(buffer_arg) {
  return service_pb.Status.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_Window(arg) {
  if (!(arg instanceof service_pb.Window)) {
    throw new Error('Expected argument of type Window');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_Window(buffer_arg) {
  return service_pb.Window.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_WindowEmitStatus(arg) {
  if (!(arg instanceof service_pb.WindowEmitStatus)) {
    throw new Error('Expected argument of type WindowEmitStatus');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_WindowEmitStatus(buffer_arg) {
  return service_pb.WindowEmitStatus.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_WindowTypeRead(arg) {
  if (!(arg instanceof service_pb.WindowTypeRead)) {
    throw new Error('Expected argument of type WindowTypeRead');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_WindowTypeRead(buffer_arg) {
  return service_pb.WindowTypeRead.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_WindowTypes(arg) {
  if (!(arg instanceof service_pb.WindowTypes)) {
    throw new Error('Expected argument of type WindowTypes');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_WindowTypes(buffer_arg) {
  return service_pb.WindowTypes.deserializeBinary(new Uint8Array(buffer_arg));
}


// OrcaCore is the central orchestration service that:
// - Manages the lifecycle of processing windows
// - Coordinates algorithm execution across distributed processors
// - Tracks DAG dependencies and execution state
// - Routes results between dependent algorithms
var OrcaCoreService = exports.OrcaCoreService = {
  // Register a processor node and its supported algorithms
registerProcessor: {
    path: '/OrcaCore/RegisterProcessor',
    requestStream: false,
    responseStream: false,
    requestType: service_pb.ProcessorRegistration,
    responseType: service_pb.Status,
    requestSerialize: serialize_ProcessorRegistration,
    requestDeserialize: deserialize_ProcessorRegistration,
    responseSerialize: serialize_Status,
    responseDeserialize: deserialize_Status,
  },
  // Submit a window for processing
emitWindow: {
    path: '/OrcaCore/EmitWindow',
    requestStream: false,
    responseStream: false,
    requestType: service_pb.Window,
    responseType: service_pb.WindowEmitStatus,
    requestSerialize: serialize_Window,
    requestDeserialize: deserialize_Window,
    responseSerialize: serialize_WindowEmitStatus,
    responseDeserialize: deserialize_WindowEmitStatus,
  },
  // Data operations
readWindowTypes: {
    path: '/OrcaCore/ReadWindowTypes',
    requestStream: false,
    responseStream: false,
    requestType: service_pb.WindowTypeRead,
    responseType: service_pb.WindowTypes,
    requestSerialize: serialize_WindowTypeRead,
    requestDeserialize: deserialize_WindowTypeRead,
    responseSerialize: serialize_WindowTypes,
    responseDeserialize: deserialize_WindowTypes,
  },
};

exports.OrcaCoreClient = grpc.makeGenericClientConstructor(OrcaCoreService, 'OrcaCore');
// Core operations
// ---------------------------- Core Operations ----------------------------
//
// OrcaProcessor defines the interface that each processing node must implement.
// Processors are language-agnostic services that:
// - Execute individual algorithms
// - Handle their own internal state
// - Report results back to the orchestrator
// Orca will schedule processors asynchronously as per the DAG
var OrcaProcessorService = exports.OrcaProcessorService = {
  // Execute part of a DAG with streaming results
// Server streams back execution results as they become available
executeDagPart: {
    path: '/OrcaProcessor/ExecuteDagPart',
    requestStream: false,
    responseStream: true,
    requestType: service_pb.ExecutionRequest,
    responseType: service_pb.ExecutionResult,
    requestSerialize: serialize_ExecutionRequest,
    requestDeserialize: deserialize_ExecutionRequest,
    responseSerialize: serialize_ExecutionResult,
    responseDeserialize: deserialize_ExecutionResult,
  },
  // Check health/status of processor. i.e. a heartbeat
healthCheck: {
    path: '/OrcaProcessor/HealthCheck',
    requestStream: false,
    responseStream: false,
    requestType: service_pb.HealthCheckRequest,
    responseType: service_pb.HealthCheckResponse,
    requestSerialize: serialize_HealthCheckRequest,
    requestDeserialize: deserialize_HealthCheckRequest,
    responseSerialize: serialize_HealthCheckResponse,
    responseDeserialize: deserialize_HealthCheckResponse,
  },
};

exports.OrcaProcessorClient = grpc.makeGenericClientConstructor(OrcaProcessorService, 'OrcaProcessor');
