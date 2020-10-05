
/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as messages_pb from "./messages_pb";

export class TelemetryMessageReceived extends jspb.Message { 
    getConfigdigest(): Uint8Array | string;
    getConfigdigest_asU8(): Uint8Array;
    getConfigdigest_asB64(): string;
    setConfigdigest(value: Uint8Array | string): TelemetryMessageReceived;


    hasMsg(): boolean;
    clearMsg(): void;
    getMsg(): messages_pb.MessageWrapper | undefined;
    setMsg(value?: messages_pb.MessageWrapper): TelemetryMessageReceived;

    getSender(): number;
    setSender(value: number): TelemetryMessageReceived;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TelemetryMessageReceived.AsObject;
    static toObject(includeInstance: boolean, msg: TelemetryMessageReceived): TelemetryMessageReceived.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TelemetryMessageReceived, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TelemetryMessageReceived;
    static deserializeBinaryFromReader(message: TelemetryMessageReceived, reader: jspb.BinaryReader): TelemetryMessageReceived;
}

export namespace TelemetryMessageReceived {
    export type AsObject = {
        configdigest: Uint8Array | string,
        msg?: messages_pb.MessageWrapper.AsObject,
        sender: number,
    }
}

export class TelemetryMessageSent extends jspb.Message { 
    getConfigdigest(): Uint8Array | string;
    getConfigdigest_asU8(): Uint8Array;
    getConfigdigest_asB64(): string;
    setConfigdigest(value: Uint8Array | string): TelemetryMessageSent;


    hasMsg(): boolean;
    clearMsg(): void;
    getMsg(): messages_pb.MessageWrapper | undefined;
    setMsg(value?: messages_pb.MessageWrapper): TelemetryMessageSent;

    getReceiver(): number;
    setReceiver(value: number): TelemetryMessageSent;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TelemetryMessageSent.AsObject;
    static toObject(includeInstance: boolean, msg: TelemetryMessageSent): TelemetryMessageSent.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TelemetryMessageSent, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TelemetryMessageSent;
    static deserializeBinaryFromReader(message: TelemetryMessageSent, reader: jspb.BinaryReader): TelemetryMessageSent;
}

export namespace TelemetryMessageSent {
    export type AsObject = {
        configdigest: Uint8Array | string,
        msg?: messages_pb.MessageWrapper.AsObject,
        receiver: number,
    }
}

export class TelemetryAssertionViolation extends jspb.Message { 

    hasInvalidsignature(): boolean;
    clearInvalidsignature(): void;
    getInvalidsignature(): TelemetryAssertionViolationInvalidSignature | undefined;
    setInvalidsignature(value?: TelemetryAssertionViolationInvalidSignature): TelemetryAssertionViolation;


    getECase(): TelemetryAssertionViolation.ECase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TelemetryAssertionViolation.AsObject;
    static toObject(includeInstance: boolean, msg: TelemetryAssertionViolation): TelemetryAssertionViolation.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TelemetryAssertionViolation, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TelemetryAssertionViolation;
    static deserializeBinaryFromReader(message: TelemetryAssertionViolation, reader: jspb.BinaryReader): TelemetryAssertionViolation;
}

export namespace TelemetryAssertionViolation {
    export type AsObject = {
        invalidsignature?: TelemetryAssertionViolationInvalidSignature.AsObject,
    }

    export enum ECase {
        E_NOT_SET = 0,
    
    INVALIDSIGNATURE = 1,

    }

}

export class TelemetryAssertionViolationInvalidSignature extends jspb.Message { 
    getConfigdigest(): Uint8Array | string;
    getConfigdigest_asU8(): Uint8Array;
    getConfigdigest_asB64(): string;
    setConfigdigest(value: Uint8Array | string): TelemetryAssertionViolationInvalidSignature;


    hasMsg(): boolean;
    clearMsg(): void;
    getMsg(): messages_pb.MessageWrapper | undefined;
    setMsg(value?: messages_pb.MessageWrapper): TelemetryAssertionViolationInvalidSignature;

    getSender(): number;
    setSender(value: number): TelemetryAssertionViolationInvalidSignature;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TelemetryAssertionViolationInvalidSignature.AsObject;
    static toObject(includeInstance: boolean, msg: TelemetryAssertionViolationInvalidSignature): TelemetryAssertionViolationInvalidSignature.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TelemetryAssertionViolationInvalidSignature, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TelemetryAssertionViolationInvalidSignature;
    static deserializeBinaryFromReader(message: TelemetryAssertionViolationInvalidSignature, reader: jspb.BinaryReader): TelemetryAssertionViolationInvalidSignature;
}

export namespace TelemetryAssertionViolationInvalidSignature {
    export type AsObject = {
        configdigest: Uint8Array | string,
        msg?: messages_pb.MessageWrapper.AsObject,
        sender: number,
    }
}

export class TelemetryStateUpdate extends jspb.Message { 
    getConfigdigest(): Uint8Array | string;
    getConfigdigest_asU8(): Uint8Array;
    getConfigdigest_asB64(): string;
    setConfigdigest(value: Uint8Array | string): TelemetryStateUpdate;

    getEpoch(): number;
    setEpoch(value: number): TelemetryStateUpdate;

    getRound(): number;
    setRound(value: number): TelemetryStateUpdate;

    getTime(): number;
    setTime(value: number): TelemetryStateUpdate;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TelemetryStateUpdate.AsObject;
    static toObject(includeInstance: boolean, msg: TelemetryStateUpdate): TelemetryStateUpdate.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TelemetryStateUpdate, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TelemetryStateUpdate;
    static deserializeBinaryFromReader(message: TelemetryStateUpdate, reader: jspb.BinaryReader): TelemetryStateUpdate;
}

export namespace TelemetryStateUpdate {
    export type AsObject = {
        configdigest: Uint8Array | string,
        epoch: number,
        round: number,
        time: number,
    }
}
