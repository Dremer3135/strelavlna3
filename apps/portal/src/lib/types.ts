import type { ContestsRecord, ContestsResponse, ConstantsRecord, ConstantsResponse, ProbsRecord, ProbsResponse } from "./pocketbase-types"

export type EditableProb = {
    prob: ProbsResponse,
    edit: Partial<ProbsRecord>
}

export type EditableConstant = {
    constant: ConstantsResponse,
    edit: Partial<ConstantsRecord>
}

export type EditableContest = {
    contest: ContestsResponse,
    edit: Partial<ContestsRecord>
}

export type LatexSegment = {
    type: 'text' | 'latex';
    content: string
};

export type DropdownItem = {
    id: string,
    value: string
}
