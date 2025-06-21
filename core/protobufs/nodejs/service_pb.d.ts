// package: 
// file: service.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as vendor_validate_pb from "./vendor/validate_pb";

export class Window extends jspb.Message {
  getTimeFrom(): number;
  setTimeFrom(value: number): void;

  getTimeTo(): number;
  setTimeTo(value: number): void;

  getWindowTypeName(): string;
  setWindowTypeName(value: string): void;

  getWindowTypeVersion(): string;
  setWindowTypeVersion(value: string): void;

  getOrigin(): string;
  setOrigin(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Window.AsObject;
  static toObject(includeInstance: boolean, msg: Window): Window.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Window, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Window;
  static deserializeBinaryFromReader(message: Window, reader: jspb.BinaryReader): Window;
}

export namespace Window {
  export type AsObject = {
    timeFrom: number,
    timeTo: number,
    windowTypeName: string,
    windowTypeVersion: string,
    origin: string,
  }
}

export class WindowType extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getVersion(): string;
  setVersion(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WindowType.AsObject;
  static toObject(includeInstance: boolean, msg: WindowType): WindowType.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WindowType, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WindowType;
  static deserializeBinaryFromReader(message: WindowType, reader: jspb.BinaryReader): WindowType;
}

export namespace WindowType {
  export type AsObject = {
    name: string,
    version: string,
  }
}

export class WindowEmitStatus extends jspb.Message {
  getStatus(): WindowEmitStatus.StatusEnumMap[keyof WindowEmitStatus.StatusEnumMap];
  setStatus(value: WindowEmitStatus.StatusEnumMap[keyof WindowEmitStatus.StatusEnumMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WindowEmitStatus.AsObject;
  static toObject(includeInstance: boolean, msg: WindowEmitStatus): WindowEmitStatus.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WindowEmitStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WindowEmitStatus;
  static deserializeBinaryFromReader(message: WindowEmitStatus, reader: jspb.BinaryReader): WindowEmitStatus;
}

export namespace WindowEmitStatus {
  export type AsObject = {
    status: WindowEmitStatus.StatusEnumMap[keyof WindowEmitStatus.StatusEnumMap],
  }

  export interface StatusEnumMap {
    NO_TRIGGERED_ALGORITHMS: 0;
    PROCESSING_TRIGGERED: 1;
    TRIGGERING_FAILED: 2;
  }

  export const StatusEnum: StatusEnumMap;
}

export class AlgorithmDependency extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getVersion(): string;
  setVersion(value: string): void;

  getProcessorName(): string;
  setProcessorName(value: string): void;

  getProcessorRuntime(): string;
  setProcessorRuntime(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AlgorithmDependency.AsObject;
  static toObject(includeInstance: boolean, msg: AlgorithmDependency): AlgorithmDependency.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AlgorithmDependency, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AlgorithmDependency;
  static deserializeBinaryFromReader(message: AlgorithmDependency, reader: jspb.BinaryReader): AlgorithmDependency;
}

export namespace AlgorithmDependency {
  export type AsObject = {
    name: string,
    version: string,
    processorName: string,
    processorRuntime: string,
  }
}

export class Algorithm extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getVersion(): string;
  setVersion(value: string): void;

  hasWindowType(): boolean;
  clearWindowType(): void;
  getWindowType(): WindowType | undefined;
  setWindowType(value?: WindowType): void;

  clearDependenciesList(): void;
  getDependenciesList(): Array<AlgorithmDependency>;
  setDependenciesList(value: Array<AlgorithmDependency>): void;
  addDependencies(value?: AlgorithmDependency, index?: number): AlgorithmDependency;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Algorithm.AsObject;
  static toObject(includeInstance: boolean, msg: Algorithm): Algorithm.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Algorithm, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Algorithm;
  static deserializeBinaryFromReader(message: Algorithm, reader: jspb.BinaryReader): Algorithm;
}

export namespace Algorithm {
  export type AsObject = {
    name: string,
    version: string,
    windowType?: WindowType.AsObject,
    dependenciesList: Array<AlgorithmDependency.AsObject>,
  }
}

export class FloatArray extends jspb.Message {
  clearValuesList(): void;
  getValuesList(): Array<number>;
  setValuesList(value: Array<number>): void;
  addValues(value: number, index?: number): number;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FloatArray.AsObject;
  static toObject(includeInstance: boolean, msg: FloatArray): FloatArray.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FloatArray, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FloatArray;
  static deserializeBinaryFromReader(message: FloatArray, reader: jspb.BinaryReader): FloatArray;
}

export namespace FloatArray {
  export type AsObject = {
    valuesList: Array<number>,
  }
}

export class Result extends jspb.Message {
  getStatus(): ResultStatusMap[keyof ResultStatusMap];
  setStatus(value: ResultStatusMap[keyof ResultStatusMap]): void;

  hasSingleValue(): boolean;
  clearSingleValue(): void;
  getSingleValue(): number;
  setSingleValue(value: number): void;

  hasFloatValues(): boolean;
  clearFloatValues(): void;
  getFloatValues(): FloatArray | undefined;
  setFloatValues(value?: FloatArray): void;

  hasStructValue(): boolean;
  clearStructValue(): void;
  getStructValue(): google_protobuf_struct_pb.Struct | undefined;
  setStructValue(value?: google_protobuf_struct_pb.Struct): void;

  getTimestamp(): number;
  setTimestamp(value: number): void;

  getResultDataCase(): Result.ResultDataCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Result.AsObject;
  static toObject(includeInstance: boolean, msg: Result): Result.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Result, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Result;
  static deserializeBinaryFromReader(message: Result, reader: jspb.BinaryReader): Result;
}

export namespace Result {
  export type AsObject = {
    status: ResultStatusMap[keyof ResultStatusMap],
    singleValue: number,
    floatValues?: FloatArray.AsObject,
    structValue?: google_protobuf_struct_pb.Struct.AsObject,
    timestamp: number,
  }

  export enum ResultDataCase {
    RESULT_DATA_NOT_SET = 0,
    SINGLE_VALUE = 2,
    FLOAT_VALUES = 3,
    STRUCT_VALUE = 4,
  }
}

export class ProcessorRegistration extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getRuntime(): string;
  setRuntime(value: string): void;

  getConnectionStr(): string;
  setConnectionStr(value: string): void;

  clearSupportedAlgorithmsList(): void;
  getSupportedAlgorithmsList(): Array<Algorithm>;
  setSupportedAlgorithmsList(value: Array<Algorithm>): void;
  addSupportedAlgorithms(value?: Algorithm, index?: number): Algorithm;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProcessorRegistration.AsObject;
  static toObject(includeInstance: boolean, msg: ProcessorRegistration): ProcessorRegistration.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProcessorRegistration, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProcessorRegistration;
  static deserializeBinaryFromReader(message: ProcessorRegistration, reader: jspb.BinaryReader): ProcessorRegistration;
}

export namespace ProcessorRegistration {
  export type AsObject = {
    name: string,
    runtime: string,
    connectionStr: string,
    supportedAlgorithmsList: Array<Algorithm.AsObject>,
  }
}

export class ProcessingTask extends jspb.Message {
  getTaskId(): string;
  setTaskId(value: string): void;

  hasAlgorithm(): boolean;
  clearAlgorithm(): void;
  getAlgorithm(): Algorithm | undefined;
  setAlgorithm(value?: Algorithm): void;

  hasWindow(): boolean;
  clearWindow(): void;
  getWindow(): Window | undefined;
  setWindow(value?: Window): void;

  clearDependencyResultsList(): void;
  getDependencyResultsList(): Array<Result>;
  setDependencyResultsList(value: Array<Result>): void;
  addDependencyResults(value?: Result, index?: number): Result;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProcessingTask.AsObject;
  static toObject(includeInstance: boolean, msg: ProcessingTask): ProcessingTask.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProcessingTask, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProcessingTask;
  static deserializeBinaryFromReader(message: ProcessingTask, reader: jspb.BinaryReader): ProcessingTask;
}

export namespace ProcessingTask {
  export type AsObject = {
    taskId: string,
    algorithm?: Algorithm.AsObject,
    window?: Window.AsObject,
    dependencyResultsList: Array<Result.AsObject>,
  }
}

export class ExecutionRequest extends jspb.Message {
  getExecId(): string;
  setExecId(value: string): void;

  hasWindow(): boolean;
  clearWindow(): void;
  getWindow(): Window | undefined;
  setWindow(value?: Window): void;

  clearAlgorithmResultsList(): void;
  getAlgorithmResultsList(): Array<AlgorithmResult>;
  setAlgorithmResultsList(value: Array<AlgorithmResult>): void;
  addAlgorithmResults(value?: AlgorithmResult, index?: number): AlgorithmResult;

  clearAlgorithmsList(): void;
  getAlgorithmsList(): Array<Algorithm>;
  setAlgorithmsList(value: Array<Algorithm>): void;
  addAlgorithms(value?: Algorithm, index?: number): Algorithm;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecutionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ExecutionRequest): ExecutionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExecutionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecutionRequest;
  static deserializeBinaryFromReader(message: ExecutionRequest, reader: jspb.BinaryReader): ExecutionRequest;
}

export namespace ExecutionRequest {
  export type AsObject = {
    execId: string,
    window?: Window.AsObject,
    algorithmResultsList: Array<AlgorithmResult.AsObject>,
    algorithmsList: Array<Algorithm.AsObject>,
  }
}

export class ExecutionResult extends jspb.Message {
  getExecId(): string;
  setExecId(value: string): void;

  hasAlgorithmResult(): boolean;
  clearAlgorithmResult(): void;
  getAlgorithmResult(): AlgorithmResult | undefined;
  setAlgorithmResult(value?: AlgorithmResult): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecutionResult.AsObject;
  static toObject(includeInstance: boolean, msg: ExecutionResult): ExecutionResult.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ExecutionResult, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecutionResult;
  static deserializeBinaryFromReader(message: ExecutionResult, reader: jspb.BinaryReader): ExecutionResult;
}

export namespace ExecutionResult {
  export type AsObject = {
    execId: string,
    algorithmResult?: AlgorithmResult.AsObject,
  }
}

export class AlgorithmResult extends jspb.Message {
  hasAlgorithm(): boolean;
  clearAlgorithm(): void;
  getAlgorithm(): Algorithm | undefined;
  setAlgorithm(value?: Algorithm): void;

  hasResult(): boolean;
  clearResult(): void;
  getResult(): Result | undefined;
  setResult(value?: Result): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AlgorithmResult.AsObject;
  static toObject(includeInstance: boolean, msg: AlgorithmResult): AlgorithmResult.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AlgorithmResult, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AlgorithmResult;
  static deserializeBinaryFromReader(message: AlgorithmResult, reader: jspb.BinaryReader): AlgorithmResult;
}

export namespace AlgorithmResult {
  export type AsObject = {
    algorithm?: Algorithm.AsObject,
    result?: Result.AsObject,
  }
}

export class Status extends jspb.Message {
  getReceived(): boolean;
  setReceived(value: boolean): void;

  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Status.AsObject;
  static toObject(includeInstance: boolean, msg: Status): Status.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Status, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Status;
  static deserializeBinaryFromReader(message: Status, reader: jspb.BinaryReader): Status;
}

export namespace Status {
  export type AsObject = {
    received: boolean,
    message: string,
  }
}

export class HealthCheckRequest extends jspb.Message {
  getTimestamp(): number;
  setTimestamp(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HealthCheckRequest.AsObject;
  static toObject(includeInstance: boolean, msg: HealthCheckRequest): HealthCheckRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HealthCheckRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HealthCheckRequest;
  static deserializeBinaryFromReader(message: HealthCheckRequest, reader: jspb.BinaryReader): HealthCheckRequest;
}

export namespace HealthCheckRequest {
  export type AsObject = {
    timestamp: number,
  }
}

export class HealthCheckResponse extends jspb.Message {
  getStatus(): HealthCheckResponse.StatusMap[keyof HealthCheckResponse.StatusMap];
  setStatus(value: HealthCheckResponse.StatusMap[keyof HealthCheckResponse.StatusMap]): void;

  getMessage(): string;
  setMessage(value: string): void;

  hasMetrics(): boolean;
  clearMetrics(): void;
  getMetrics(): ProcessorMetrics | undefined;
  setMetrics(value?: ProcessorMetrics): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HealthCheckResponse.AsObject;
  static toObject(includeInstance: boolean, msg: HealthCheckResponse): HealthCheckResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HealthCheckResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HealthCheckResponse;
  static deserializeBinaryFromReader(message: HealthCheckResponse, reader: jspb.BinaryReader): HealthCheckResponse;
}

export namespace HealthCheckResponse {
  export type AsObject = {
    status: HealthCheckResponse.StatusMap[keyof HealthCheckResponse.StatusMap],
    message: string,
    metrics?: ProcessorMetrics.AsObject,
  }

  export interface StatusMap {
    STATUS_UNKNOWN: 0;
    STATUS_SERVING: 1;
    STATUS_TRANSITIONING: 2;
    STATUS_NOT_SERVING: 3;
  }

  export const Status: StatusMap;
}

export class ProcessorMetrics extends jspb.Message {
  getActiveTasks(): number;
  setActiveTasks(value: number): void;

  getMemoryBytes(): number;
  setMemoryBytes(value: number): void;

  getCpuPercent(): number;
  setCpuPercent(value: number): void;

  getUptimeSeconds(): number;
  setUptimeSeconds(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProcessorMetrics.AsObject;
  static toObject(includeInstance: boolean, msg: ProcessorMetrics): ProcessorMetrics.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProcessorMetrics, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProcessorMetrics;
  static deserializeBinaryFromReader(message: ProcessorMetrics, reader: jspb.BinaryReader): ProcessorMetrics;
}

export namespace ProcessorMetrics {
  export type AsObject = {
    activeTasks: number,
    memoryBytes: number,
    cpuPercent: number,
    uptimeSeconds: number,
  }
}

export class WindowTypeRead extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WindowTypeRead.AsObject;
  static toObject(includeInstance: boolean, msg: WindowTypeRead): WindowTypeRead.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WindowTypeRead, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WindowTypeRead;
  static deserializeBinaryFromReader(message: WindowTypeRead, reader: jspb.BinaryReader): WindowTypeRead;
}

export namespace WindowTypeRead {
  export type AsObject = {
  }
}

export class WindowTypes extends jspb.Message {
  clearWindowsList(): void;
  getWindowsList(): Array<WindowType>;
  setWindowsList(value: Array<WindowType>): void;
  addWindows(value?: WindowType, index?: number): WindowType;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WindowTypes.AsObject;
  static toObject(includeInstance: boolean, msg: WindowTypes): WindowTypes.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WindowTypes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WindowTypes;
  static deserializeBinaryFromReader(message: WindowTypes, reader: jspb.BinaryReader): WindowTypes;
}

export namespace WindowTypes {
  export type AsObject = {
    windowsList: Array<WindowType.AsObject>,
  }
}

export interface ResultStatusMap {
  RESULT_STATUS_HANDLED_FAILED: 0;
  RESULT_STATUS_UNHANDLED_FAILED: 1;
  RESULT_STATUS_SUCEEDED: 2;
}

export const ResultStatus: ResultStatusMap;

