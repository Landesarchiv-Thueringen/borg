import { CommonModule, PercentPipe } from '@angular/common';
import { Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { BreakOpportunitiesPipe } from '../pipes/break-opportunities.pipe';
import { FileFeaturePipe } from '../pipes/file-feature.pipe';
import { FeatureSet, FeatureValue, ToolResult } from '../results';

interface DialogData {
  featureSets: FeatureSet[];
  toolResults: ToolResult[];
}

interface FeatureWithMarks {
  value: string | boolean | number;
  toolMarks: ToolMark[];
}

interface Row {
  puid: FeatureWithMarks | undefined;
  mimeType: FeatureWithMarks | undefined;
  version: FeatureWithMarks | undefined;
  valid: FeatureWithMarks | undefined;
  score: number;
  tools: ToolMark[];
}

interface ToolLegendEntry {
  id: string;
  title: string;
  mark: ToolMark;
}

interface ToolMark {
  value: string;
  class: string;
}

@Component({
  selector: 'app-feature-sets-table',
  templateUrl: './feature-sets-table.component.html',
  styleUrls: ['./feature-sets-table.component.scss'],
  imports: [
    CommonModule,
    MatDialogModule,
    MatButtonModule,
    FileFeaturePipe,
    MatTableModule,
    PercentPipe,
    BreakOpportunitiesPipe,
    MatIconModule,
  ],
})
export class FeatureSetsTableComponent {
  private readonly data = inject<DialogData>(MAT_DIALOG_DATA);
  displayedColumns: string[] = ['tools', 'puid', 'mimeType', 'version', 'valid', 'score'];
  dataSource: Row[] = [];
  toolLegend: ToolLegendEntry[] = [];
  constructor() {
    for (let [index, tr] of this.data.toolResults.entries()) {
      const mark: ToolMark = {
        value: tr.title.charAt(0),
        class: 'tool-mark' + (index + 1),
      };
      this.toolLegend.push({
        id: tr.id,
        title: tr.title,
        mark: mark,
      });
    }
    this.dataSource = this.data.featureSets.map((fs) => {
      return {
        puid: this.getFeatureWithMarks(fs, 'format:puid'),
        mimeType: this.getFeatureWithMarks(fs, 'format:mimeType'),
        version: this.getFeatureWithMarks(fs, 'format:version'),
        valid: this.getFeatureWithMarks(fs, 'format:valid'),
        score: fs.score,
        tools: this.getFeatureSetMarks(fs),
      };
    });
  }

  getFeatureSetMarks(fs: FeatureSet): ToolMark[] {
    return fs.supportingTools.map((toolId) => this.getToolMark(toolId));
  }

  getFeatureWithMarks(fs: FeatureSet, key: string): FeatureWithMarks | undefined {
    const f = fs.features[key];
    if (!f) {
      return undefined;
    }
    return {
      value: f.value,
      toolMarks: this.getToolMarks(f),
    };
  }

  getToolMarks(f: FeatureValue): ToolMark[] {
    return f.supportingTools.map((toolId) => {
      return this.getToolMark(toolId);
    });
  }

  getToolMark(toolId: string): ToolMark {
    const entry = this.toolLegend.find((tool) => tool.id === toolId);
    if (entry) {
      return entry.mark;
    }
    const mark: ToolMark = {
      value: '[' + (toolId + 1) + ']',
      class: '',
    };
    return mark;
  }
}
