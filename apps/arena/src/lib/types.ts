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
    start: Date,
    end: Date,
    runningState: "before" | "running" | "after"
    rank: number,
}

export type MessageType = {
    origin: "recieved" | "sent",
    type: "message" | "answer" | "grade" | "window-focus" | "copy" | "paste",
    value: string,
    sentTime: Date
}
