import type { Prob } from "./types";

export function isProbSolved(prob: Prob): boolean {
    return prob.chat.find(p => p.type == "grade" && p.value == "correct") !== undefined;
}
export function isProbUngraded(prob: Prob): boolean {
    return prob.chat.findLastIndex(p => p.type == "answer") > prob.chat.findLastIndex(p => p.type == "grade");
}
export function getProbLastAnswer(prob: Prob): string | undefined {
    return prob.chat.findLast(p => p.type == "answer")?.value ?? undefined;
}
export function hasProbChat(prob: Prob): boolean {
    return prob.chat.length > 0;
}