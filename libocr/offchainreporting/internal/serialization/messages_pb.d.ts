
/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class MessageNewEpoch extends jspb.Message { 
    getEpoch(): number;
    setEpoch(value: number): MessageNewEpoch;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MessageNewEpoch.AsObject;
    static toObject(includeInstance: boolean, msg: MessageNewEpoch): MessageNewEpoch.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MessageNewEpoch, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MessageNewEpoch;
    static deserializeBinaryFromReader(message: MessageNewEpoch, reader: jspb.BinaryReader): MessageNewEpoch;
}

export namespace MessageNewEpoch {
    export type AsObject = {
        epoch: number,
    }
}

export class MessageObserveReq extends jspb.Message { 
    getRound(): number;
    setRound(value: number): MessageObserveReq;

    getEpoch(): number;
    setEpoch(value: number): MessageObserveReq;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MessageObserveReq.AsObject;
    static toObject(includeInstance: boolean, msg: MessageObserveReq): MessageObserveReq.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MessageObserveReq, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MessageObserveReq;
    static deserializeBinaryFromReader(message: MessageObserveReq, reader: jspb.BinaryReader): MessageObserveReq;
}

export namespace MessageObserveReq {
    export type AsObject = {
        round: number,
        epoch: number,
    }
}

export class ReportingContext extends jspb.Message { 
    getConfigdigest(): Uint8Array | string;
    getConfigdigest_asU8(): Uint8Array;
    getConfigdigest_asB64(): string;
    setConfigdigest(value: Uint8Array | string): ReportingContext;

    getEpoch(): number;
    setEpoch(value: number): ReportingContext;

    getRound(): number;
    setRound(value: number): ReportingContext;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ReportingContext.AsObject;
    static toObject(includeInstance: boolean, msg: ReportingContext): ReportingContext.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ReportingContext, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ReportingContext;
    static deserializeBinaryFromReader(message: ReportingContext, reader: jspb.BinaryReader): ReportingContext;
}

export namespace ReportingContext {
    export type AsObject = {
        configdigest: Uint8Array | string,
        epoch: number,
        round: number,
    }
}

export class ObservationValue extends jspb.Message { 
    getValue(): Uint8Array | string;
    getValue_asU8(): Uint8Array;
    getValue_asB64(): string;
    setValue(value: Uint8Array | string): ObservationValue;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ObservationValue.AsObject;
    static toObject(includeInstance: boolean, msg: ObservationValue): ObservationValue.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ObservationValue, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ObservationValue;
    static deserializeBinaryFromReader(message: ObservationValue, reader: jspb.BinaryReader): ObservationValue;
}

export namespace ObservationValue {
    export type AsObject = {
        value: Uint8Array | string,
    }
}

export class Observation extends jspb.Message { 

    hasCtx(): boolean;
    clearCtx(): void;
    getCtx(): ReportingContext | undefined;
    setCtx(value?: ReportingContext): Observation;


    hasValue(): boolean;
    clearValue(): void;
    getValue(): ObservationValue | undefined;
    setValue(value?: ObservationValue): Observation;

    getSignature(): Uint8Array | string;
    getSignature_asU8(): Uint8Array;
    getSignature_asB64(): string;
    setSignature(value: Uint8Array | string): Observation;

    getOracleid(): number;
    setOracleid(value: number): Observation;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Observation.AsObject;
    static toObject(includeInstance: boolean, msg: Observation): Observation.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Observation, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Observation;
    static deserializeBinaryFromReader(message: Observation, reader: jspb.BinaryReader): Observation;
}

export namespace Observation {
    export type AsObject = {
        ctx?: ReportingContext.AsObject,
        value?: ObservationValue.AsObject,
        signature: Uint8Array | string,
        oracleid: number,
    }
}

export class MessageObserve extends jspb.Message { 
    getEpoch(): number;
    setEpoch(value: number): MessageObserve;

    getRound(): number;
    setRound(value: number): MessageObserve;


    hasObs(): boolean;
    clearObs(): void;
    getObs(): Observation | undefined;
    setObs(value?: Observation): MessageObserve;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MessageObserve.AsObject;
    static toObject(includeInstance: boolean, msg: MessageObserve): MessageObserve.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MessageObserve, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MessageObserve;
    static deserializeBinaryFromReader(message: MessageObserve, reader: jspb.BinaryReader): MessageObserve;
}

export namespace MessageObserve {
    export type AsObject = {
        epoch: number,
        round: number,
        obs?: Observation.AsObject,
    }
}

export class MessageReportReq extends jspb.Message { 
    getRound(): number;
    setRound(value: number): MessageReportReq;

    getEpoch(): number;
    setEpoch(value: number): MessageReportReq;

    clearObservationsList(): void;
    getObservationsList(): Array<Observation>;
    setObservationsList(value: Array<Observation>): MessageReportReq;
    addObservations(value?: Observation, index?: number): Observation;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MessageReportReq.AsObject;
    static toObject(includeInstance: boolean, msg: MessageReportReq): MessageReportReq.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MessageReportReq, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MessageReportReq;
    static deserializeBinaryFromReader(message: MessageReportReq, reader: jspb.BinaryReader): MessageReportReq;
}

export namespace MessageReportReq {
    export type AsObject = {
        round: number,
        epoch: number,
        observationsList: Array<Observation.AsObject>,
    }
}

export class Signature extends jspb.Message { 
    getSignature(): Uint8Array | string;
    getSignature_asU8(): Uint8Array;
    getSignature_asB64(): string;
    setSignature(value: Uint8Array | string): Signature;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Signature.AsObject;
    static toObject(includeInstance: boolean, msg: Signature): Signature.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Signature, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Signature;
    static deserializeBinaryFromReader(message: Signature, reader: jspb.BinaryReader): Signature;
}

export namespace Signature {
    export type AsObject = {
        signature: Uint8Array | string,
    }
}

export class OracleValue extends jspb.Message { 
    getOracleid(): number;
    setOracleid(value: number): OracleValue;


    hasValue(): boolean;
    clearValue(): void;
    getValue(): ObservationValue | undefined;
    setValue(value?: ObservationValue): OracleValue;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): OracleValue.AsObject;
    static toObject(includeInstance: boolean, msg: OracleValue): OracleValue.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: OracleValue, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): OracleValue;
    static deserializeBinaryFromReader(message: OracleValue, reader: jspb.BinaryReader): OracleValue;
}

export namespace OracleValue {
    export type AsObject = {
        oracleid: number,
        value?: ObservationValue.AsObject,
    }
}

export class ContractReport extends jspb.Message { 

    hasCtx(): boolean;
    clearCtx(): void;
    getCtx(): ReportingContext | undefined;
    setCtx(value?: ReportingContext): ContractReport;

    clearValuesList(): void;
    getValuesList(): Array<OracleValue>;
    setValuesList(value: Array<OracleValue>): ContractReport;
    addValues(value?: OracleValue, index?: number): OracleValue;

    getSig(): Uint8Array | string;
    getSig_asU8(): Uint8Array;
    getSig_asB64(): string;
    setSig(value: Uint8Array | string): ContractReport;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ContractReport.AsObject;
    static toObject(includeInstance: boolean, msg: ContractReport): ContractReport.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ContractReport, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ContractReport;
    static deserializeBinaryFromReader(message: ContractReport, reader: jspb.BinaryReader): ContractReport;
}

export namespace ContractReport {
    export type AsObject = {
        ctx?: ReportingContext.AsObject,
        valuesList: Array<OracleValue.AsObject>,
        sig: Uint8Array | string,
    }
}

export class MessageReport extends jspb.Message { 
    getEpoch(): number;
    setEpoch(value: number): MessageReport;

    getRound(): number;
    setRound(value: number): MessageReport;


    hasContractreport(): boolean;
    clearContractreport(): void;
    getContractreport(): ContractReport | undefined;
    setContractreport(value?: ContractReport): MessageReport;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MessageReport.AsObject;
    static toObject(includeInstance: boolean, msg: MessageReport): MessageReport.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MessageReport, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MessageReport;
    static deserializeBinaryFromReader(message: MessageReport, reader: jspb.BinaryReader): MessageReport;
}

export namespace MessageReport {
    export type AsObject = {
        epoch: number,
        round: number,
        contractreport?: ContractReport.AsObject,
    }
}

export class ContractReportWithSignatures extends jspb.Message { 

    hasContractreport(): boolean;
    clearContractreport(): void;
    getContractreport(): ContractReport | undefined;
    setContractreport(value?: ContractReport): ContractReportWithSignatures;

    clearSignaturesList(): void;
    getSignaturesList(): Array<Signature>;
    setSignaturesList(value: Array<Signature>): ContractReportWithSignatures;
    addSignatures(value?: Signature, index?: number): Signature;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ContractReportWithSignatures.AsObject;
    static toObject(includeInstance: boolean, msg: ContractReportWithSignatures): ContractReportWithSignatures.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ContractReportWithSignatures, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ContractReportWithSignatures;
    static deserializeBinaryFromReader(message: ContractReportWithSignatures, reader: jspb.BinaryReader): ContractReportWithSignatures;
}

export namespace ContractReportWithSignatures {
    export type AsObject = {
        contractreport?: ContractReport.AsObject,
        signaturesList: Array<Signature.AsObject>,
    }
}

export class MessageFinal extends jspb.Message { 
    getEpoch(): number;
    setEpoch(value: number): MessageFinal;

    getLeader(): number;
    setLeader(value: number): MessageFinal;

    getRound(): number;
    setRound(value: number): MessageFinal;


    hasReport(): boolean;
    clearReport(): void;
    getReport(): ContractReportWithSignatures | undefined;
    setReport(value?: ContractReportWithSignatures): MessageFinal;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MessageFinal.AsObject;
    static toObject(includeInstance: boolean, msg: MessageFinal): MessageFinal.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MessageFinal, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MessageFinal;
    static deserializeBinaryFromReader(message: MessageFinal, reader: jspb.BinaryReader): MessageFinal;
}

export namespace MessageFinal {
    export type AsObject = {
        epoch: number,
        leader: number,
        round: number,
        report?: ContractReportWithSignatures.AsObject,
    }
}

export class MessageFinalEcho extends jspb.Message { 

    hasFinal(): boolean;
    clearFinal(): void;
    getFinal(): MessageFinal | undefined;
    setFinal(value?: MessageFinal): MessageFinalEcho;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MessageFinalEcho.AsObject;
    static toObject(includeInstance: boolean, msg: MessageFinalEcho): MessageFinalEcho.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MessageFinalEcho, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MessageFinalEcho;
    static deserializeBinaryFromReader(message: MessageFinalEcho, reader: jspb.BinaryReader): MessageFinalEcho;
}

export namespace MessageFinalEcho {
    export type AsObject = {
        pb_final?: MessageFinal.AsObject,
    }
}

export class MessageWrapper extends jspb.Message { 

    hasMessagenewepoch(): boolean;
    clearMessagenewepoch(): void;
    getMessagenewepoch(): MessageNewEpoch | undefined;
    setMessagenewepoch(value?: MessageNewEpoch): MessageWrapper;


    hasMessageobservereq(): boolean;
    clearMessageobservereq(): void;
    getMessageobservereq(): MessageObserveReq | undefined;
    setMessageobservereq(value?: MessageObserveReq): MessageWrapper;


    hasMessageobserve(): boolean;
    clearMessageobserve(): void;
    getMessageobserve(): MessageObserve | undefined;
    setMessageobserve(value?: MessageObserve): MessageWrapper;


    hasMessagereportreq(): boolean;
    clearMessagereportreq(): void;
    getMessagereportreq(): MessageReportReq | undefined;
    setMessagereportreq(value?: MessageReportReq): MessageWrapper;


    hasMessagereport(): boolean;
    clearMessagereport(): void;
    getMessagereport(): MessageReport | undefined;
    setMessagereport(value?: MessageReport): MessageWrapper;


    hasMessagefinal(): boolean;
    clearMessagefinal(): void;
    getMessagefinal(): MessageFinal | undefined;
    setMessagefinal(value?: MessageFinal): MessageWrapper;


    hasMessagefinalecho(): boolean;
    clearMessagefinalecho(): void;
    getMessagefinalecho(): MessageFinalEcho | undefined;
    setMessagefinalecho(value?: MessageFinalEcho): MessageWrapper;


    getMsgCase(): MessageWrapper.MsgCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MessageWrapper.AsObject;
    static toObject(includeInstance: boolean, msg: MessageWrapper): MessageWrapper.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: MessageWrapper, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MessageWrapper;
    static deserializeBinaryFromReader(message: MessageWrapper, reader: jspb.BinaryReader): MessageWrapper;
}

export namespace MessageWrapper {
    export type AsObject = {
        messagenewepoch?: MessageNewEpoch.AsObject,
        messageobservereq?: MessageObserveReq.AsObject,
        messageobserve?: MessageObserve.AsObject,
        messagereportreq?: MessageReportReq.AsObject,
        messagereport?: MessageReport.AsObject,
        messagefinal?: MessageFinal.AsObject,
        messagefinalecho?: MessageFinalEcho.AsObject,
    }

    export enum MsgCase {
        MSG_NOT_SET = 0,
    
    MESSAGENEWEPOCH = 2,

    MESSAGEOBSERVEREQ = 3,

    MESSAGEOBSERVE = 4,

    MESSAGEREPORTREQ = 5,

    MESSAGEREPORT = 6,

    MESSAGEFINAL = 7,

    MESSAGEFINALECHO = 8,

    }

}
