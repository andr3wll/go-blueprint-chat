package main

import (
	"errors"
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
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar struct
type AuthAvatar struct{}

// UseAuthAvatar variable
var UseAuthAvatar AuthAvatar

// GetAvatarURL method for type AuthAvatar
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar_url"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	urlStr, ok := url.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}

	return urlStr, nil

	// if url, ok := c.userData["avatar_url"]; ok {
	// 	if urlStr, ok := url.(string); ok {
	// 		return urlStr, nil
	// 	}
	// }
	// return "", ErrNoAvatarURL
}

// GravatarAvatar struct
type GravatarAvatar struct{}

// UseGravatarAvatar variable
var UseGravatarAvatar GravatarAvatar

// GetAvatarURL method for type GravatarAvatar
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// FileSystemAvatar struct
type FileSystemAvatar struct{}

// UseFileSystemAvatar variable
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL method for type FileSystemAvatar
func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "/avatars/" + useridStr + ".jpg", nil
		}
	}
	return "", ErrNoAvatarURL
}
