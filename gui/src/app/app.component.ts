import { Component } from '@angular/core';
import { Router, RouterOutlet } from '@angular/router';
import { MainNavigationComponent } from './main-navigation/main-navigation.component';
import { FileDropContainerComponent } from './utility/file-drop-container/file-drop-container.component';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  standalone: true,
  imports: [RouterOutlet, MainNavigationComponent, FileDropContainerComponent],
})
export class AppComponent {
  constructor(router: Router) {
    // Any results are lost when the app (re-)loads. Navigate to the upload
    // page, which is the only useful page now.
    router.navigate(['.']);
  }
}
