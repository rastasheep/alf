import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

// containers
import { EditorComponent } from './containers/editor/editor.component';

// components
// import { TeslaCarComponent } from './components/tesla-car/tesla-car.component';

// services
// import { BatteryService } from './tesla-battery.service';

export const ROUTES: Routes = [{
  path: 'editor',
  children: [
    { path: '', redirectTo: 'schema', pathMatch: 'full' },
    { path: 'schema', component: EditorComponent, data: { section: 'schema' }},
    { path: 'history', component: EditorComponent, data: { section: 'history' }},
    { path: 'saved', component: EditorComponent, data: { section: 'saved' }},
    { path: 'saved/:id', component: EditorComponent, data: { section: 'saved' }}
  ]
}];

@NgModule({
  declarations: [
    EditorComponent,
    // TeslaCarComponent,
  ],
  imports: [
    CommonModule,
    RouterModule.forRoot(ROUTES),
  ],
  providers: [
    // BatteryService
  ],
  exports: [
    EditorComponent
  ]
})
export class EditorModule {}
