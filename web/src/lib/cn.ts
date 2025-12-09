import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import type { ClassValue } from 'clsx';

/**
 * Utility function to merge class names into a single string.
 * Uses clsx for conditional classes and tailwind-merge to handle conflicts.
 * @param inputs - The class names to merge.
 * @returns The merged class names.
 */
export function cn(...inputs: Array<ClassValue>) {
  return twMerge(clsx(inputs));
}
