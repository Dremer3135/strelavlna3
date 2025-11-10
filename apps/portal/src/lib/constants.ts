// src/lib/constants.ts

import { ProbsDiffOptions } from "./pocketbase-types";

/**
 * This file contains shared constants and inferred types that are used
 * across the application. Defining them here provides a single source of truth.
 */

// --- Problem Difficulties ---

// The runtime array of possible difficulty values, derived directly from the PocketBase schema.
export const PROB_DIFFICULTIES = Object.values(ProbsDiffOptions);

// The compile-time TypeScript type, also derived from the PocketBase schema.
export type ProbDifficulty = ProbsDiffOptions;


// --- Add other shared constants below ---

// Example for user roles:
/*
export const USER_ROLES = ['admin', 'corrector', 'student'] as const;
export type UserRole = (typeof USER_ROLES)[number];
*/
