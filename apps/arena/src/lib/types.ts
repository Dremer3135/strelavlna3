import type { ProbContentType } from "../../../shared/dist/types";
export type ContestState = "waiting" | "running" | "ended";

export type Prob = ProbContentType & {
    focusedBy: string[],
    chat: MessageType[]
};

export type CurrentState = {
    teamName: string,
    money: number,
    myId: string
}

export type MessageType = {
    origin: "recieved" | "sent",
    type: "message" | "answer" | "grade" | "window-focus" | "copy" | "paste",
    value: string
}