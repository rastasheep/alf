import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { DashboardComponent } from './containers/dashboard/dashboard.component';

// components
// import { TeslaCarComponent } from './components/tesla-car/tesla-car.component';

// services
// import { BatteryService } from './tesla-battery.service';

export const ROUTES: Routes = [{
  path: 'editor',
  children: [
    { path: '', redirectTo: 'schema', pathMatch: 'full' },
    { path: 'schema', component: DashboardComponent, data: { section: 'schema' }},
    { path: 'history', component: DashboardComponent, data: { section: 'history' }},
    { path: 'saved', component: DashboardComponent, data: { section: 'saved' }},
    { path: 'saved/:id', component: DashboardComponent, data: { section: 'saved' }}
  ]
}];

@NgModule({
  declarations: [
    DashboardComponent,
    // TeslaCarComponent,
  ],
  imports: [
    CommonModule,
    RouterModule.forRoot(ROUTES),
  ],
  providers: [
    // BatteryService
  ],
})
export class DashboardModule {}
