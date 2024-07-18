export interface FileResult {
  id: string;
  fileName: string;
  relativePath?: string;
  fileSize: number;
  toolResults: ToolResults;
}

export interface ToolResults {
  fileIdentificationResults: ToolResult[] | null;
  fileValidationResults: ToolResult[] | null;
  summary: Summary;
}

export interface ToolResult {
  toolName: string;
  toolVersion: string;
  toolOutput: string;
  outputFormat: 'text' | 'json' | 'csv';
  extractedFeatures: { [key: string]: string };
  error: string | null;
}

export interface Summary {
  [key: string]: Feature;
}

export interface Feature {
  key: string;
  values: FeatureValue[];
}

export interface FeatureValue {
  value: string;
  score: number;
  tools: ToolConfidence[];
}

export interface ToolConfidence {
  confidence: number;
  toolName: string;
}
