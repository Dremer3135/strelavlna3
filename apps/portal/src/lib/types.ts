import type { ProbsRecord, ProbsResponse } from "./pocketbase-types"

export type EditableProb = {
    prob: ProbsResponse,
    edit: Partial<ProbsRecord>
}