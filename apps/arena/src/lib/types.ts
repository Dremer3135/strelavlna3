import type { ProbContentType } from "../../../shared/dist/types";
export type ContestState = "waiting" | "running" | "ended";

export type Prob = ProbContentType & {
    focusedBy: string[],
    chat: MessageType[],
    owned: "bought" | "solved" | "sold",
};

export type CurrentState = {
    teamName: string,
    money: number,
    myId: string,
    probsRemaining: number[],
    pricesBuy: number[],
    pricesSell: number[],
    procesSolve: number[],
    start: Date | undefined,
    end: Date | undefined,
    runningState: "before" | "running" | "after" | "results" | "paused"
    rank: number,
    isAdmin?: boolean
    results?: Record<string, ResultsAtom>,
}

export type MessageType = {
    origin: "recieved" | "sent",
    type: "message" | "answer" | "grade" | "window-focus" | "copy" | "paste",
    value: string,
    sentTime: Date
}

export type ResultsAtom = {
  teamName: string,
  rank: number,
  money: number,
}
