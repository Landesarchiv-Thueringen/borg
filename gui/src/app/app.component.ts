import { Component, inject } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { FileDropContainerComponent } from './core/file-drop-container/file-drop-container.component';
import { MainNavigationComponent } from './core/main-navigation/main-navigation.component';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss'],
    imports: [RouterOutlet, MainNavigationComponent, FileDropContainerComponent]
})
export class AppComponent {
  constructor() {
    const router = inject(Router);

    // Any results are lost when the app (re-)loads. Navigate to the upload
    // page, which is the only useful page now.
    router.navigate(['.']);
  }
}
