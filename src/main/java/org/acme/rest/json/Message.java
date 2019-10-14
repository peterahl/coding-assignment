package org.acme.rest.json;

import java.util.Objects;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
public class Message {

  private long id;

  private String text;

  public Message() {
  }

  public Message(int id, String text) {
    this.id = id;
    this.text = text;
  }

  @Id
  @GeneratedValue(strategy = GenerationType.SEQUENCE, generator="messageSeq")
  public long getId() {
    return id;
  }

  public void setId(long id) {
    this.id = id;
  }

  public String getText() {
    return text;
  }

  public void setText(String text) {
    this.text = text;
  }

  @Override
  public boolean equals(Object obj) {
    if (!(obj instanceof Message)) {
      return false;
    }

    Message other = (Message) obj;

    return Objects.equals(other.id, this.id);
  }

  @Override
  public int hashCode() {
    return Objects.hash(this.id);
  }
}
