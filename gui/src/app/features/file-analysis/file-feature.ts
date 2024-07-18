const OVERVIEW_FEATURES = [
  'relativePath',
  'fileName',
  'fileSize',
  'puid',
  'mimeType',
  'formatVersion',
  'valid',
] as const;
export type OverviewFeature = (typeof OVERVIEW_FEATURES)[number];

const featureOrder = new Map<string, number>([
  ['relativePath', 1],
  ['fileName', 2],
  ['fileSize', 3],
  ['puid', 4],
  ['mimeType', 5],
  ['formatVersion', 6],
  ['encoding', 7],
  ['', 101],
  ['wellFormed', 1001],
  ['valid', 1002],
]);

export function isOverviewFeature(feature: string): feature is OverviewFeature {
  return (OVERVIEW_FEATURES as readonly string[]).includes(feature);
}

/** Sorts feature keys and removes duplicates. */
export function sortFeatures(features: string[]): string[] {
  features = [...new Set(features)];
  return features.sort((f1: string, f2: string) => {
    let orderF1: number | undefined = featureOrder.get(f1);
    if (!orderF1) {
      orderF1 = featureOrder.get('');
    }
    let orderF2: number | undefined = featureOrder.get(f2);
    if (!orderF2) {
      orderF2 = featureOrder.get('');
    }
    if (orderF1! < orderF2!) {
      return -1;
    } else if (orderF1! > orderF2!) {
      return 1;
    }
    return 0;
  });
}
