import type { ConstsRecord, ConstsResponse, ProbsRecord, ProbsResponse } from "./pocketbase-types"

export type EditableProb = {
    prob: ProbsResponse,
    edit: Partial<ProbsRecord>
}

export type EditableConst = {
    const: ConstsResponse,
    edit: Partial<ConstsRecord>
}