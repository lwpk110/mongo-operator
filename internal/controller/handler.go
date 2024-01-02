package controller

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CheckAndCreateNs(ctx context.Context, client client.Client, nameSpace string) error {
	/*	ns := &corev1.Namespace{}
		client.Get(ctx)*/
	return nil
}
