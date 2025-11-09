import type { ConstantsRecord, ConstantsResponse, ProbsRecord, ProbsResponse } from "./pocketbase-types"

export type EditableProb = {
    prob: ProbsResponse,
    edit: Partial<ProbsRecord>
}

export type EditableConstant = {
    constant: ConstantsResponse,
    edit: Partial<ConstantsRecord>
}

export type LatexSegment = {
    type: 'text' | 'latex';
    content: string;
};