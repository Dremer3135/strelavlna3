import { writable, derived } from 'svelte/store';
import type { CurrentState, Prob } from '$lib/types';
import { probs } from './probs';

export const currentState = writable<CurrentState>({
    teamName: "-",
    money: 0,
    myId: "bablba",
    probsRemaining: [10, 2, 0],
    pricesBuy: [10, 30, 80],
    pricesSell: [10, 15, 40],
    procesSolve: [15, 50, 200],
    start: undefined,
    end: undefined,
    runningState: "before",
    rank: 0,
    isAdmin: false
});

export const focusedProb = derived<[typeof probs, typeof currentState], Prob | undefined>(
    [probs, currentState],
    ([$probs, $currentState]) => {
        return Object.values($probs).find(prob => prob.focusedBy.includes($currentState.myId))
    }
);

export const wsConnected = writable<boolean>(true);
