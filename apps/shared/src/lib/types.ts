export type ProbContentType = {
    name: string,
    text: string,
    answer: string | undefined,
    images: string[],
    diff: string,
    id: string
}


export type LatexSegment = {
    type: 'text' | 'latex';
    content: string
};
