/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"strconv"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/google/go-github/github"
	"github.com/yquansah/github-operator/api/v1alpha1"
	vcsv1alpha1 "github.com/yquansah/github-operator/api/v1alpha1"
)

// GitRepositoryReconciler reconciles a GitRepository object
type GitRepositoryReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	GithubClient *github.Client
}

const githubRepoFinalizer = "github-repo-finalizer"

//+kubebuilder:rbac:groups=vcs.github,resources=gitrepositories,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=vcs.github,resources=gitrepositories/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=vcs.github,resources=gitrepositories/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GitRepository object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile

// +kubebuilder:rbac:groups=vcs.github,resources=gitrepositories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vcs.github,resources=gitrepositories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=vcs.github,resources=gitrepositories/finalizers,verbs=update
func (r *GitRepositoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	ghRepo := &v1alpha1.GitRepository{}

	err := r.Get(ctx, req.NamespacedName, ghRepo)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	isGitHubRepoToBeDeleted := ghRepo.GetDeletionTimestamp() != nil
	if isGitHubRepoToBeDeleted {
		if controllerutil.ContainsFinalizer(ghRepo, githubRepoFinalizer) {
			err = r.finalizeGithubRepo(ctx, ghRepo)
			if err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(ghRepo, githubRepoFinalizer)
			err := r.Update(ctx, ghRepo)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	if !ghRepo.Status.Created {
		repo, _, err := r.GithubClient.Repositories.Create(ctx, "", &github.Repository{
			Name:        &ghRepo.Spec.Name,
			Description: &ghRepo.Spec.Description,
			Private:     &ghRepo.Spec.Private,
		})
		if err != nil {
			return ctrl.Result{}, err
		}

		ghRepo.Status.Created = true
		ghRepo.Status.URL = repo.GetURL()
		ghRepo.Status.ID = strconv.Itoa(int(repo.GetID()))

		err = r.Status().Update(ctx, ghRepo)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if !controllerutil.ContainsFinalizer(ghRepo, githubRepoFinalizer) {
		controllerutil.AddFinalizer(ghRepo, githubRepoFinalizer)
		err = r.Update(ctx, ghRepo)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *GitRepositoryReconciler) finalizeGithubRepo(ctx context.Context, ghRepo *v1alpha1.GitRepository) error {
	_, err := r.GithubClient.Repositories.Delete(ctx, "yquansah", ghRepo.Spec.Name)

	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *GitRepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vcsv1alpha1.GitRepository{}).
		Complete(r)
}
