package shell

import (
	"yunion.io/x/onecloud/pkg/util/azure"
	"yunion.io/x/onecloud/pkg/util/shellutils"
)

func init() {
	type ImageListOptions struct {
		Limit  int `help:"page size"`
		Offset int `help:"page offset"`
	}
	shellutils.R(&ImageListOptions{}, "image-list", "List images", func(cli *azure.SRegion, args *ImageListOptions) error {
		if images, err := cli.GetImages(); err != nil {
			return err
		} else {
			printList(images, len(images), args.Offset, args.Limit, []string{})
			return nil
		}
	})

	type ImageCreateOptions struct {
		NAME       string `helo:"Image name"`
		OSTYPE     string `helo:"Operation system" choices:"Linux|Windows"`
		BLOBURI    string `helo:"page blob uri"`
		DISKSIZEGB int32  `helo:"Image size"`
	}

	shellutils.R(&ImageCreateOptions{}, "image-create", "Create image", func(cli *azure.SRegion, args *ImageCreateOptions) error {
		if image, err := cli.CreateImageByBlob(args.NAME, args.OSTYPE, args.BLOBURI, args.DISKSIZEGB); err != nil {
			return err
		} else {
			printObject(image)
			return nil
		}
	})

	type ImageDeleteOptions struct {
		ID string `helo:"Image ID"`
	}

	shellutils.R(&ImageDeleteOptions{}, "image-delete", "Delete image", func(cli *azure.SRegion, args *ImageDeleteOptions) error {
		return cli.DeleteImage(args.ID)
	})

}
