import type { ConstantsResponse, CorrectorsResponse, ProbsResponse, TeachersResponse } from '$lib/types/pocketbase-types';
import type { EditableConstant, EditableProb, LatexSegment } from './types';

type AppUser = CorrectorsResponse | TeachersResponse;

/**
 * Type Guard to check if a user is a Corrector.
 * @param user The user object to check.
 * @returns True if the user is a CorrectorsResponse, false otherwise.
 */
export const isCorrector = (user: AppUser | undefined): user is CorrectorsResponse => {
    return user?.collectionName === 'correctors';
}


export const isProbEdited = (eprob: EditableProb): boolean => {
    return Object.keys(eprob.edit).length > 0
}

export const isConstantEdited = (econstant: EditableConstant): boolean => {
    return Object.keys(econstant.edit).length > 0
}

export const getProbEditedState = (eprob: EditableProb): ProbsResponse => {
    return {
        ...eprob.prob,
        ...(eprob.edit as Partial<ProbsResponse>)
    };
};

export const getConstantEditedState = (econstant: EditableConstant): ConstantsResponse => {
    return {
        ...econstant.constant,
        ...(econstant.edit as Partial<ConstantsResponse>)
    };
};

export function filterRecord<T>( 
    record: Record<string, T>,
    predicate: (id: string, value: T) => boolean
): Record<string, T> {
    return Object.entries(record)
        .filter(([id, value]) => predicate(id, value as T))
        .reduce((acc, [id, value]) => {
            acc[id] = value as T;
            return acc;
        }, {} as Record<string, T>);
}


export function parseLatex(text: string): LatexSegment[] {
    const segments: LatexSegment[] = [];
    const regex = /\$(.*?)\$/g;
    let lastIndex = 0;
    let match;

    while ((match = regex.exec(text)) !== null) {
        // Push the text before the match
        if (match.index > lastIndex) {
            segments.push({ type: 'text', content: text.substring(lastIndex, match.index) });
        }
        // Push the latex part, using the captured group
        segments.push({ type: 'latex', content: match[1] });
        lastIndex = regex.lastIndex;
    }

    // Push any remaining text after the last match
    if (lastIndex < text.length) {
        segments.push({ type: 'text', content: text.substring(lastIndex) });
    }

    return segments;
}