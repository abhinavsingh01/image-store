import { Component } from '@angular/core';
import { ApiService } from 'src/app/services/api.service';
import { Router, ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent {
  username: string = '';
  password: any = '';
  errorMessage: any = '';

  constructor(private api: ApiService, private router: Router) {}

  ngOnInit() {
    if (localStorage.getItem('access_token') != undefined) {
      this.router.navigate(['dashboard']);
    }
  }

  register(): void {
    this.router.navigate(['register']);
  }

  login(): void {
    this.api.loginCall(this.username, this.password).subscribe(
      (response) => {
        const data = response.data;
        const token = data['token'];
        debugger;
        localStorage.setItem('access_token', token);
        this.router.navigate(['dashboard']);
      },
      (error) => {
        this.errorMessage = error;
      }
    );
  }
}
