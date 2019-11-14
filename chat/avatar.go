package main

import (
	"errors"
	"io/ioutil"
	"path"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(ChatUser) (string, error)
}

// TryAvatars represents slice of Avatar objects
type TryAvatars []Avatar

// GetAvatarURL method for type TryAvatars
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

// AuthAvatar struct
type AuthAvatar struct{}

// UseAuthAvatar variable
var UseAuthAvatar AuthAvatar

// GetAvatarURL method for type AuthAvatar
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

// GravatarAvatar struct
type GravatarAvatar struct{}

// UseGravatar variable
var UseGravatar GravatarAvatar

// GetAvatarURL method for type GravatarAvatar
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// FileSystemAvatar struct
type FileSystemAvatar struct{}

// UseFileSystemAvatar variable
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL method for type FileSystemAvatar
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
			return "/avatars/" + file.Name(), nil
		}

	}
	return "", ErrNoAvatarURL
}
